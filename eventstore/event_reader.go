package eventstore

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"time"

	uuid "github.com/satori/go.uuid"
)

type EventReader struct {
	registry     map[string]Initialiser
	file         *os.File
	scanner      *bufio.Scanner
	currentIndex int
}

type EventMeta struct {
	Timestamp       time.Time `json:"meta_timestamp"`
	ID              uuid.UUID `json:"meta_id"`
	AggregateRootID uuid.UUID `json:"meta_aggregate_id"`
	Version         int       `json:"meta_version"`
	Type            string    `json:"meta_type"`
}

func (e *EventMeta) event() *EventMeta      { return e }
func (e *EventMeta) AggregateID() uuid.UUID { return e.AggregateRootID }

type Event interface {
	event() *EventMeta
	AggregateID() uuid.UUID
}

func NewEventReader(registry map[string]Initialiser, filename string) (*EventReader, error) {
	file, err := os.Open(filename)

	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	scanner := bufio.NewScanner(file)

	return &EventReader{registry, file, scanner, 0}, nil
}

func (er *EventReader) Close() error {
	return er.file.Close()
}

func (er *EventReader) ReadAll() bool {
	er.currentIndex++
	return er.scanner.Scan()
}

func (er *EventReader) ReadFor(uuid uuid.UUID) bool {

	for er.scanner.Scan() {
		er.currentIndex++
		record, err := er.Event()

		if err != nil {
			return false
		}

		if record.event().AggregateRootID == uuid {
			return true
		}
	}

	return false
}

func (er *EventReader) ReadFrom(offset int) bool {

	if er.currentIndex < offset {
		for i := 0; i < offset; i++ {
			er.scanner.Scan()
		}
		er.currentIndex = offset
	}

	return er.scanner.Scan()
}

type eventType struct {
	Type string `json:"meta_type"`
}

func (er *EventReader) Event() (Event, error) {
	var et eventType
	if err := json.Unmarshal(er.scanner.Bytes(), &et); err != nil {
		return nil, err
	}

	creator, found := er.registry[et.Type]

	if !found {
		return nil, fmt.Errorf("Unable to find an event of type %s", et.Type)
	}

	event := creator()
	if err := json.Unmarshal(er.scanner.Bytes(), event); err != nil {
		return nil, err
	}

	return event.(Event), nil
}
