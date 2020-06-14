package eventstore

import (
	"reflect"

	uuid "github.com/satori/go.uuid"
)

type Projector func(state interface{}, event Event) interface{}

// ------------

type EventStore interface {
	RegisterEvent(creator Initialiser)
	RegisterProjection(name string, initialiseState Initialiser, project Projector)
	ReadView(name string, view interface{}) error
	LoadAggregate(id uuid.UUID, a Aggregate) error
	SaveAggregate(a Aggregate) error
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

func (es *eventStore) RegisterEvent(creator Initialiser) {
	es.registry[EventName(creator())] = creator
}

func (es *eventStore) RegisterProjection(name string, initialiseState Initialiser, project Projector) {
	es.projections[name] = projection{
		name:            name,
		initialiseState: initialiseState,
		projector:       project,
	}
}

func (es *eventStore) ReadView(name string, view interface{}) error {
	v := es.backend.NewView(name)

	return v.ReadView(view)
}

func (es *eventStore) LoadAggregate(id uuid.UUID, a Aggregate) error {

	er, err := es.backend.NewEventReader(es.registry)
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

func (es *eventStore) SaveAggregate(a Aggregate) error {

	writer := es.backend.NewEventWriter()
	aggregate := a.aggregator()

	newVersion, err := writer.WriteEvents(aggregate.id, aggregate.version, aggregate.changes)

	if err != nil {
		return err
	}

	aggregate.changes = []Event{}
	aggregate.version = newVersion

	return es.runProjections()
}

func (es *eventStore) runProjections() error {

	views := es.allViews()
	lowestIndex, err := findUnprocessedEvents(views)

	if err != nil {
		return err
	}

	events, err := es.loadEvents(lowestIndex)
	if err != nil {
		return err
	}

	lastIndex := lowestIndex + len(events)

	for name, projection := range es.projections {

		view := views[name]

		// we will have already failed if this didn't work
		viewLastIndex, _ := view.LastEventIndex()
		projectionEvents := events[viewLastIndex-lowestIndex:]

		if len(projectionEvents) == 0 {
			continue
		}

		state := projection.initialiseState()
		if err := view.ReadView(state); err != nil {
			return err
		}

		for _, e := range projectionEvents {
			state = projection.projector(state, e)
		}

		return view.WriteView(state, lastIndex)
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

func findUnprocessedEvents(views map[string]View) (int, error) {

	lowestIndex := 0

	for _, view := range views {
		index, err := view.LastEventIndex()
		if err != nil {
			return 0, err
		}

		lowestIndex = min(lowestIndex, index)
	}

	return lowestIndex, nil
}

func (es *eventStore) loadEvents(lowestIndex int) ([]Event, error) {
	er, err := es.backend.NewEventReader(es.registry)
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
