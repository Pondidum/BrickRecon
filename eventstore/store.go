package eventstore

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"reflect"
	"time"

	uuid "github.com/satori/go.uuid"
)

var newline = []byte("\n")

type EventStore struct {
	root string

	registry    map[string]func() interface{}
	projections map[string]Projector
}

type Event struct {
	Timestamp time.Time
	ID        string
	Type      string
	Content   interface{}
}

type Projector func(e interface{}) interface{}

type readEvent struct {
	Timestamp time.Time
	ID        string
	Type      string
	Content   json.RawMessage
}

func CreateEventStore(root string) *EventStore {
	return &EventStore{
		root:        root,
		registry:    map[string]func() interface{}{},
		projections: map[string]Projector{},
	}
}

func (es *EventStore) RegisterEvent(creator func() interface{}) {
	es.registry[eventName(creator())] = creator
}

func eventName(event interface{}) string {
	t := reflect.TypeOf(event)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	return t.Name()
}

func (es *EventStore) Write(events ...interface{}) error {

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

		for name, p := range es.projections {
			view := p(e)
			viewBytes, err := json.Marshal(view)

			if err != nil {
				return err
			}

			err = os.MkdirAll(path.Join(es.root, "views"), os.ModePerm)
			if err != nil {
				return err
			}

			err = ioutil.WriteFile(path.Join(es.root, "views", name+".json"), viewBytes, 0666)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (es *EventStore) ReadEvents(offset int) ([]interface{}, error) {

	filename := path.Join(es.root, "events")
	file, err := os.Open(filename)

	if err != nil {
		return nil, err
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	events := []interface{}{}

	for scanner.Scan() {

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

func (es *EventStore) RegisterProjection(name string, projection Projector) {
	es.projections[name] = projection
}
