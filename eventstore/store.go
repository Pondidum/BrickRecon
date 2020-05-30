package eventstore

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"reflect"
	"strconv"
	"time"

	uuid "github.com/satori/go.uuid"
)

var newline = []byte("\n")

type EventStore struct {
	root string

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
		registry:    map[string]Initialiser{},
		projections: map[string]Projection{},
	}
}

func (es *EventStore) RegisterEvent(creator Initialiser) {
	es.registry[eventName(creator())] = creator
}

func (es *EventStore) RegisterProjection(name string, initialisState Initialiser, project Projector) {
	es.projections[name] = NewProjection(path.Join(es.root, "views"), name, initialisState, project)
}

func eventName(event interface{}) string {
	t := reflect.TypeOf(event)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	return t.Name()
}

func (es *EventStore) checkIndex(filename string) (int, error) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return 0, err
	}

	return strconv.Atoi(string(content))
}

func (es *EventStore) updateCheckIndex(filename string, index int) error {
	return ioutil.WriteFile(filename, []byte(strconv.Itoa(index)), 0666)
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

	checkIndexPath := path.Join(es.root, "checkindex")
	checkindex, err := es.checkIndex(checkIndexPath)

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

		checkindex++
	}

	if err = es.updateCheckIndex(checkIndexPath, checkindex); err != nil {
		return err
	}

	return es.runProjections()
}

func (es *EventStore) runProjections() error {

	lowestIndex := 0
	projectionIndex := map[string]int{}

	for name, projection := range es.projections {

		if index, err := projection.CheckIndex(); err != nil {
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

	for _, projection := range es.projections {

		firstEvent := projectionIndex[projection.name] - lowestIndex

		if err := projection.Project(events[firstEvent:]); err != nil {
			return err
		}

		if err := projection.WriteCheckIndex(lastIndex); err != nil {
			return err
		}
	}

	return nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
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
