package eventstore

import (
	"bufio"
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

type Event struct {
	Timestamp time.Time
	ID        string
	Type      string
	Content   interface{}
}

type Initialiser func() interface{}

type readEvent struct {
	Timestamp time.Time
	ID        string
	Type      string
	Content   json.RawMessage
}

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

func (es *EventStore) Write(event interface{}) error {
	return es.WriteEvents([]interface{}{event})
}

func (es *EventStore) WriteEvents(events []interface{}) error {

	filename := path.Join(es.root, "events")
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)

	if err != nil {
		return err
	}

	defer file.Close()

	for _, e := range events {

		dto := &Event{
			Timestamp: time.Now(),
			ID:        uuid.NewV4().String(),
			Type:      eventName(e),
			Content:   e,
		}

		bytes, err := json.Marshal(dto)

		if err != nil {
			return err
		}

		if _, err := file.Write(append(bytes, newline...)); err != nil {
			return err
		}
	}

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

	events, err := es.ReadEvents(lowestIndex)
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

func (es *EventStore) ReadView(name string, view interface{}) error {
	p := es.projections[name]
	return p.ReadView(view)
}

func (es *EventStore) ReadEvents(offset int) ([]interface{}, error) {

	filename := path.Join(es.root, "events")
	file, err := os.Open(filename)

	if err != nil {
		return nil, err
	}

	defer file.Close()

	events := []interface{}{}
	scanner := bufio.NewScanner(file)
	lines := 0

	for scanner.Scan() {

		if lines < offset {
			lines++
			continue
		}

		var read readEvent
		if err := json.Unmarshal(scanner.Bytes(), &read); err != nil {
			return nil, err
		}

		creator, found := es.registry[read.Type]

		if !found {
			return nil, fmt.Errorf("Unable to find an event of type %s", read.Type)
		}

		e := creator()

		if err := json.Unmarshal(read.Content, &e); err != nil {
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
