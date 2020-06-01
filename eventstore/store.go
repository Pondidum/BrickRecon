package eventstore

import (
	"bytes"
	"encoding/json"
	"os"
	"path"
	"reflect"
	"time"

	uuid "github.com/satori/go.uuid"
)

var newline = []byte("\n")

type EventStore struct {
	root string

	checkIndex  CheckIndex
	registry    map[string]Initialiser
	projections map[string]Projection
}

type Event struct {
	ID          uuid.UUID
	Timestamp   time.Time
	AggregateID uuid.UUID
	Version     int
	Type        string
	Content     interface{}
}

type Initialiser func() interface{}

func NewEventStore(root string) *EventStore {
	return &EventStore{
		root:        root,
		checkIndex:  NewCheckIndex(),
		registry:    map[string]Initialiser{},
		projections: map[string]Projection{},
	}
}

func (es *EventStore) RegisterEvent(creator Initialiser) {
	es.registry[eventName(creator())] = creator
}

func (es *EventStore) RegisterProjection(name string, initialiseState Initialiser, project Projector) {
	es.projections[name] = NewProjection(path.Join(es.root, "views"), name, initialiseState, project)
}

func (es *EventStore) LoadAggregate(id uuid.UUID, a *Aggregator) error {

	events, err := es.readAggregateEvents(id)

	if err != nil {
		return err
	}

	a.fromEvents(events)

	return nil
}

func (es *EventStore) ReadView(name string, view interface{}) error {
	p := es.projections[name]
	return p.ReadView(view)
}

func (es *EventStore) SaveAggregate(a *Aggregator) error {

	filename := path.Join(es.root, "events")
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)

	if err != nil {
		return err
	}

	defer file.Close()

	currentVersion := a.version

	block := bytes.Buffer{}

	for _, e := range a.changes {

		currentVersion++
		dto := &Event{
			ID:          uuid.NewV4(),
			Timestamp:   time.Now(),
			AggregateID: a.id,
			Version:     currentVersion,
			Type:        eventName(e),
			Content:     e,
		}

		bytes, err := json.Marshal(dto)

		if err != nil {
			return err
		}

		if _, err := block.Write(bytes); err != nil {
			return err
		}

		if _, err := block.Write(newline); err != nil {
			return err
		}
	}

	if _, err := file.Write(block.Bytes()); err != nil {
		return err
	}

	a.changes = []interface{}{}
	a.version = currentVersion

	return es.runProjections()
}

func (es *EventStore) runProjections() error {

	if err := os.MkdirAll(path.Join(es.root, "views"), os.ModePerm); err != nil {
		return err
	}

	lowestIndex := 0
	projectionIndex := map[string]int{}

	for name, projection := range es.projections {

		if index, err := es.checkIndex.Read(projection.path); err != nil {
			return err
		} else {
			lowestIndex = min(lowestIndex, index)
			projectionIndex[name] = index
		}
	}

	events, err := es.readEvents(lowestIndex)
	if err != nil {
		return err
	}

	lastIndex := lowestIndex + len(events)

	for name, projection := range es.projections {

		firstEvent := projectionIndex[name] - lowestIndex

		if err := projection.Project(events[firstEvent:]); err != nil {
			return err
		}

		if err := es.checkIndex.Write(projection.path, lastIndex); err != nil {
			return err
		}
	}

	return nil
}

func (es *EventStore) readEvents(offset int) ([]interface{}, error) {

	er, err := NewEventReader(es.registry, path.Join(es.root, "events"))
	if err != nil {
		return nil, err
	}

	defer er.Close()

	events := []interface{}{}

	for er.ReadFrom(offset) {

		e, err := er.Event()
		if err != nil {
			return nil, err
		}

		events = append(events, e)
	}

	return events, nil
}
func (es *EventStore) readAggregateEvents(aggregateID uuid.UUID) ([]interface{}, error) {

	er, err := NewEventReader(es.registry, path.Join(es.root, "events"))
	if err != nil {
		return nil, err
	}

	defer er.Close()

	events := []interface{}{}

	for er.ReadFor(aggregateID) {

		e, err := er.Event()
		if err != nil {
			return nil, err
		}

		events = append(events, e)
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
