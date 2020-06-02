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

type Event struct {
	Timestamp   time.Time
	ID          uuid.UUID
	AggregateID uuid.UUID
	Version     int
	Type        string
}

func (e *Event) event() *Event              { return e }
func (e *Event) AggregateRootID() uuid.UUID { return e.ID }

type IsEvent interface {
	event() *Event
	AggregateRootID() uuid.UUID
}

func NewEventReader(registry map[string]Initialiser, filename string) (*EventReader, error) {
	file, err := os.Open(filename)

	if err != nil {
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
		record, err := er.Record()

		if err != nil {
			return false
		}

		if record.event().AggregateID == uuid {
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
	Type string
}

func (er *EventReader) Record() (IsEvent, error) {
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

	return event.(IsEvent), nil
}
