package eventstore

import (
	"context"
	"reflect"
	"time"

	"github.com/honeycombio/beeline-go"
	"github.com/honeycombio/beeline-go/timer"
	uuid "github.com/satori/go.uuid"
)

type Projector func(state interface{}, event Event) interface{}

// ------------

type EventStore interface {
	RegisterEvent(ctx context.Context, creator Initialiser) error
	RegisterEvents(ctx context.Context, creator []Initialiser) error
	RegisterEventMiddleware(ctx context.Context, middleware EventMiddleware) error

	RegisterProjection(ctx context.Context, projection Projection)
	ReadView(ctx context.Context, name string, view interface{}) error
	RebuildProjections(ctx context.Context) error

	LoadAggregate(ctx context.Context, id AggregateID, a Aggregate) error
	SaveAggregate(ctx context.Context, a Aggregate) error
}

type eventStore struct {
	registry    *EventRegistry
	middleware  []EventMiddleware
	projections map[string]Projection

	backend Backend
}

type Initialiser func() interface{}

func NewEventStore(backend Backend) EventStore {
	return &eventStore{
		registry:    NewRegistry(),
		middleware:  []EventMiddleware{},
		projections: map[string]Projection{},
		backend:     backend,
	}
}

type Projection interface {
	Name() string
	CreateState() interface{}
	Project(state interface{}, event Event) interface{}
}

func (es *eventStore) RegisterEvent(ctx context.Context, creator Initialiser) error {
	return es.registry.Register(ctx, creator)
}

func (es *eventStore) RegisterEvents(ctx context.Context, creators []Initialiser) error {
	return es.registry.RegisterMany(ctx, creators)
}

func (es *eventStore) RegisterEventMiddleware(ctx context.Context, middleware EventMiddleware) error {
	es.middleware = append(es.middleware, middleware)
	return nil
}

func (es *eventStore) RegisterProjection(ctx context.Context, projection Projection) {
	es.projections[projection.Name()] = projection
}

func (es *eventStore) ReadView(ctx context.Context, name string, view interface{}) error {
	var err error
	ctx, fn := buildSpan(ctx, "read_view")
	defer func() {
		fn(err)
	}()

	v := es.backend.NewView(name)
	err = v.ReadView(ctx, view)

	return err
}

func (es *eventStore) LoadAggregate(ctx context.Context, id AggregateID, a Aggregate) error {
	var err error
	ctx, fn := buildSpan(ctx, "load_aggregate")
	defer func() {
		fn(err)
	}()

	beeline.AddField(ctx, "es.aggregate_id", id)

	er, err := es.backend.NewEventReader(es.registry, id)
	if err != nil {
		return err
	}

	defer er.Close()

	aggregator := a.aggregator()
	hasEvents := false

	for er.Read() {
		hasEvents = true

		e, err := er.Event()
		if err != nil {
			return err
		}

		e = es.applyMiddleware(ctx, e)

		aggregator.onEvent(e)
		aggregator.sequence = e.Meta().Sequence
	}

	beeline.AddField(ctx, "es.aggregate_sequence", aggregator.sequence)

	if !hasEvents {
		beeline.AddField(ctx, "es.aggregate_not_found", true)
		err = &AggregateNotFoundError{ID: id}
	}

	return err
}

func (es *eventStore) SaveAggregate(ctx context.Context, a Aggregate) error {
	var err error
	ctx, fn := buildSpan(ctx, "save_aggregate")
	defer func() {
		fn(err)
	}()

	writer := es.backend.NewEventWriter()
	aggregate := a.aggregator()
	events := aggregate.changes

	beeline.AddField(ctx, "es.aggregate_id", aggregate.id)
	beeline.AddField(ctx, "es.aggregate_old_sequence", aggregate.sequence)
	beeline.AddField(ctx, "es.aggregate_changes", len(events))

	currentSequence := aggregate.sequence

	for _, e := range events {

		currentSequence++

		meta := e.Meta()

		meta.Timestamp = time.Now()
		meta.ID = uuid.NewV4()
		meta.AggregateRootID = aggregate.id
		meta.Sequence = currentSequence
		meta.Type = eventName(e)

	}

	eventsWritten, err := writer.WriteEvents(ctx, aggregate.id, events)
	if err != nil {
		return err
	}

	aggregate.changes = []Event{}
	aggregate.sequence = aggregate.sequence + eventsWritten

	beeline.AddField(ctx, "es.aggregate_sequence", aggregate.sequence)

	err = es.runProjections(ctx, events)

	return err
}

func (es *eventStore) runProjections(ctx context.Context, events []Event) error {

	var err error
	for name, projection := range es.projections {

		view := es.backend.NewView(name)
		state := projection.CreateState()
		if err := view.ReadView(ctx, state); err != nil {
			return err
		}

		for _, e := range events {
			state = projection.Project(state, e)
		}

		err = view.WriteView(ctx, state, 0)
	}

	return err
}

func (es *eventStore) RebuildProjections(ctx context.Context) error {

	be := es.backend
	if err := be.DestroyViews(); err != nil {
		return err
	}

	aggregates, err := be.AllAggregates()
	if err != nil {
		beeline.AddField(ctx, "es.err_reading_aggregates", err)
		return err
	}

	beeline.AddField(ctx, "es.aggregate_count", len(aggregates))

	for _, id := range aggregates {

		if err := es.processAggregateProjections(ctx, id); err != nil {
			return err
		}
	}

	return nil
}

func (es *eventStore) processAggregateProjections(ctx context.Context, id AggregateID) error {
	var err error
	ctx, fn := buildSpan(ctx, "process_aggregate_"+string(id))
	defer func() {
		fn(err)
	}()

	beeline.AddField(ctx, "es.aggregate_id", id)

	reader, err := es.backend.NewEventReader(es.registry, id)

	if err != nil {
		return err
	}

	defer reader.Close()

	events := []Event{}

	for reader.Read() {
		e, err := reader.Event()
		if err != nil {
			return err
		}

		e = es.applyMiddleware(ctx, e)

		events = append(events, e)
	}

	beeline.AddField(ctx, "es.aggregate_events", len(events))

	if err := es.runProjections(ctx, events); err != nil {
		return err
	}

	return nil
}

func (es *eventStore) applyMiddleware(ctx context.Context, event Event) Event {
	for _, mw := range es.middleware {
		event = mw(ctx, event)
	}

	return event
}

func eventName(event interface{}) string {
	t := reflect.TypeOf(event)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	return t.Name()
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func buildSpan(ctx context.Context, name string) (context.Context, func(error)) {
	time := timer.Start()
	c, s := beeline.StartSpan(ctx, name)

	fn := func(err error) {
		duration := time.Finish()
		if err != nil {
			beeline.AddField(c, "es.error", err.Error())
		}
		beeline.AddField(c, "es.duration_ms", duration)
		s.Send()
	}

	return c, fn
}
