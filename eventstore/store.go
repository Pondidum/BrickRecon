package eventstore

import (
	"encoding/json"
	"os"
	"path"
	"time"

	uuid "github.com/satori/go.uuid"
)

var newline = []byte("\n")

type EventStore struct {
	root string
}

type Event struct {
	Timestamp time.Time
	ID        string
	Content   interface{}
}

func CreateEventStore(root string) *EventStore {
	return &EventStore{root: root}
}

func (es *EventStore) Write(events ...interface{}) error {

	filename := path.Join(es.root, "events")
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)

	if err != nil {
		return err
	}

	defer file.Close()

	for _, e := range events {

		dto := Event{
			Timestamp: time.Now(),
			ID:        uuid.NewV4().String(),
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

	return nil
}
