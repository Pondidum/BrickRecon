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

// ------------

var newline = []byte("\n")

type EventStore struct {
	root string

	registry    map[string]Initialiser
	projections map[string]Projection

	backend Backend
}

type Initialiser func() interface{}

func NewEventStore(root string) *EventStore {
	return &EventStore{
		root:        root,
		registry:    map[string]Initialiser{},
		projections: map[string]Projection{},

		backend: &FsBackend{
			root: root,
		},
	}
}

func (es *EventStore) RegisterEvent(creator Initialiser) {
	es.registry[eventName(creator())] = creator
}

func (es *EventStore) RegisterProjection(name string, initialiseState Initialiser, project Projector) {
	es.projections[name] = NewProjection(path.Join(es.root, "views"), name, initialiseState, project)
}

func (es *EventStore) ReadView(name string, view interface{}) error {
	p := es.projections[name]
	return p.ReadView(view)
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

	lowestIndex := 0
	projectionIndex := map[string]int{}

	for name, projection := range es.projections {

		if index, err := projection.LastEventIndex(); err != nil {
			return err
		} else {
			lowestIndex = min(lowestIndex, index)
			projectionIndex[name] = index
		}
	}

	er, err := es.backend.NewEventReader(es.registry)
	if err != nil {
		return err
	}

	defer er.Close()

	events := []Event{}

	for er.ReadFrom(lowestIndex) {

		record, err := er.Event()
		if err != nil {
			return err
		}

		events = append(events, record)
	}

	for name, projection := range es.projections {

		firstEvent := projectionIndex[name] - lowestIndex

		if err := projection.Project(events[firstEvent:]); err != nil {
			return err
		}
	}

	return nil
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
