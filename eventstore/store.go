package eventstore

import (
	"context"
	"reflect"

	"github.com/honeycombio/beeline-go"
	"github.com/honeycombio/beeline-go/timer"
	uuid "github.com/satori/go.uuid"
)

type Projector func(state interface{}, event Event) interface{}

// ------------

type EventStore interface {
	RegisterEvent(ctx context.Context, creator Initialiser)
	RegisterProjection(ctx context.Context, name string, initialiseState Initialiser, project Projector)
	ReadView(ctx context.Context, name string, view interface{}) error
	LoadAggregate(ctx context.Context, id uuid.UUID, a Aggregate) error
	SaveAggregate(ctx context.Context, a Aggregate) error
}

type eventStore struct {
	registry    map[string]Initialiser
	projections map[string]projection

	backend Backend
}

type Initialiser func() interface{}

func NewEventStore(backend Backend) EventStore {
	return &eventStore{
		registry:    map[string]Initialiser{},
		projections: map[string]projection{},
		backend:     backend,
	}
}

type projection struct {
	name            string
	initialiseState Initialiser
	projector       Projector
}

func (es *eventStore) RegisterEvent(ctx context.Context, creator Initialiser) {
	es.registry[EventName(creator())] = creator
}

func (es *eventStore) RegisterProjection(ctx context.Context, name string, initialiseState Initialiser, project Projector) {
	es.projections[name] = projection{
		name:            name,
		initialiseState: initialiseState,
		projector:       project,
	}
}

func (es *eventStore) ReadView(ctx context.Context, name string, view interface{}) error {
	var err error
	ctx, fn := buildSpan(ctx)
	defer func() {
		fn(err)
	}()

	v := es.backend.NewView(name)
	err = v.ReadView(ctx, view)

	return err
}

func (es *eventStore) LoadAggregate(ctx context.Context, id uuid.UUID, a Aggregate) error {

	er, err := es.backend.NewEventReader(es.registry, ctx)
	if err != nil {
		return err
	}

	defer er.Close()

	aggregator := a.aggregator()
	hasEvents := false

	for er.ReadFor(id) {
		hasEvents = true

		r, err := er.Event()
		if err != nil {
			return err
		}

		aggregator.onEvent(r)
		aggregator.version = r.Meta().Version
	}

	if !hasEvents {
		return &AggregateNotFoundError{ID: id}
	}

	return nil
}

func (es *eventStore) SaveAggregate(ctx context.Context, a Aggregate) error {

	writer := es.backend.NewEventWriter()
	aggregate := a.aggregator()

	newVersion, err := writer.WriteEvents(ctx, aggregate.id, aggregate.version, aggregate.changes)

	if err != nil {
		return err
	}

	aggregate.changes = []Event{}
	aggregate.version = newVersion

	return es.runProjections(ctx)
}

func (es *eventStore) runProjections(ctx context.Context) error {

	views := es.allViews()
	lowestIndex, err := findUnprocessedEvents(ctx, views)

	if err != nil {
		return err
	}

	events, err := es.loadEvents(ctx, lowestIndex)
	if err != nil {
		return err
	}

	lastIndex := lowestIndex + len(events)

	for name, projection := range es.projections {

		view := views[name]

		// we will have already failed if this didn't work
		viewLastIndex, _ := view.LastEventIndex(ctx)
		projectionEvents := events[viewLastIndex-lowestIndex:]

		if len(projectionEvents) == 0 {
			continue
		}

		state := projection.initialiseState()
		if err := view.ReadView(ctx, state); err != nil {
			return err
		}

		for _, e := range projectionEvents {
			state = projection.projector(state, e)
		}

		return view.WriteView(ctx, state, lastIndex)
	}

	return nil
}

func (es *eventStore) allViews() map[string]View {
	views := make(map[string]View, len(es.projections))

	for name := range es.projections {
		views[name] = es.backend.NewView(name)
	}

	return views
}

func findUnprocessedEvents(ctx context.Context, views map[string]View) (int, error) {

	lowestIndex := 0

	for _, view := range views {
		index, err := view.LastEventIndex(ctx)
		if err != nil {
			return 0, err
		}

		lowestIndex = min(lowestIndex, index)
	}

	return lowestIndex, nil
}

func (es *eventStore) loadEvents(ctx context.Context, lowestIndex int) ([]Event, error) {
	er, err := es.backend.NewEventReader(es.registry, ctx)
	if err != nil {
		return nil, err
	}
	defer er.Close()

	events := []Event{}

	for er.ReadFrom(lowestIndex) {

		record, err := er.Event()
		if err != nil {
			return nil, err
		}

		events = append(events, record)
	}

	return events, nil
}

func EventName(event interface{}) string {
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

func buildSpan(ctx context.Context) (context.Context, func(error)) {
	time := timer.Start()
	c, s := beeline.StartSpan(ctx, "")

	fn := func(err error) {
		duration := time.Finish()
		if err != nil {
			beeline.AddField(ctx, "db.error", err.Error())
		}
		beeline.AddField(c, "es.duration_ms", duration)
		s.Send()
	}

	return c, fn
}
