package eventstore

import (
	"bytes"
	"encoding/json"
	"fmt"
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

func (es *EventStore) LoadAggregate(id uuid.UUID, a Aggregate) error {

	er, err := NewEventReader(es.registry, path.Join(es.root, "events"))
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
		return fmt.Errorf("No aggregate found for ID %s", id)
	}

	return nil
}

func (es *EventStore) ReadView(name string, view interface{}) error {
	p := es.projections[name]
	return p.ReadView(view)
}

func (es *EventStore) SaveAggregate(a Aggregate) error {

	filename := path.Join(es.root, "events")
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)

	if err != nil {
		return err
	}

	defer file.Close()

	aggregator := a.aggregator()
	currentVersion := aggregator.version

	block := bytes.Buffer{}

	for _, e := range aggregator.changes {

		currentVersion++

		meta := e.event()

		meta.Timestamp = time.Now()
		meta.ID = uuid.NewV4()
		meta.AggregateRootID = aggregator.id
		meta.Version = currentVersion
		meta.Type = eventName(e)

		bytes, err := json.Marshal(e)

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

	aggregator.changes = []Event{}
	aggregator.version = currentVersion

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

	er, err := NewEventReader(es.registry, path.Join(es.root, "events"))
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
