package eventstore

import (
	"os"
	"path"
	"reflect"

	uuid "github.com/satori/go.uuid"
)

type Backend interface {
	NewEventReader(map[string]Initialiser) (EventReader, error)
	NewEventWriter() EventWriter
	NewView(name string) View
}

type FsBackend struct {
	root string
}

func (be *FsBackend) NewEventReader(registry map[string]Initialiser) (EventReader, error) {
	return NewEventReader(registry, path.Join(be.root, "events"))
}

func (be *FsBackend) NewEventWriter() EventWriter {
	return NewEventWriter(path.Join(be.root, "events"))
}

func (be *FsBackend) NewView(name string) View {
	return &FsView{
		filename: path.Join(be.root, "views", name+".json"),
	}
}

type Projector func(state interface{}, event Event) interface{}

// ------------

var newline = []byte("\n")

type EventStore struct {
	root string

	registry    map[string]Initialiser
	projections map[string]projection

	backend Backend
}

type Initialiser func() interface{}

func NewEventStore(root string) *EventStore {
	return &EventStore{
		root:        root,
		registry:    map[string]Initialiser{},
		projections: map[string]projection{},

		backend: &FsBackend{
			root: root,
		},
	}
}

type projection struct {
	name            string
	initialiseState Initialiser
	projector       Projector
}

func (es *EventStore) RegisterEvent(creator Initialiser) {
	es.registry[eventName(creator())] = creator
}

func (es *EventStore) RegisterProjection(name string, initialiseState Initialiser, project Projector) {
	es.projections[name] = projection{
		name:            name,
		initialiseState: initialiseState,
		projector:       project,
	}
}

func (es *EventStore) ReadView(name string, view interface{}) error {
	v := es.backend.NewView(name)

	return v.ReadView(view)
}

func (es *EventStore) LoadAggregate(id uuid.UUID, a Aggregate) error {

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
		aggregator.version = r.event().Version
	}

	if !hasEvents {
		return &AggregateNotFoundError{ID: id}
	}

	return nil
}

func (es *EventStore) SaveAggregate(a Aggregate) error {

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

func (es *EventStore) runProjections() error {

	if err := os.MkdirAll(path.Join(es.root, "views"), os.ModePerm); err != nil {
		return err
	}

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

func (es *EventStore) allViews() map[string]View {
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

func (es *EventStore) loadEvents(lowestIndex int) ([]Event, error) {
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
