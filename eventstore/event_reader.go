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

type Record struct {
	Timestamp   time.Time
	ID          uuid.UUID
	AggregateID uuid.UUID
	Version     int
	Type        string
	Content     json.RawMessage

	creator func() (interface{}, error)
}

func (r *Record) Event() (interface{}, error) {
	return r.creator()
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

		if record.AggregateID == uuid {
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

func (er *EventReader) Record() (Record, error) {
	var read Record
	if err := json.Unmarshal(er.scanner.Bytes(), &read); err != nil {
		return read, err
	}

	creator, found := er.registry[read.Type]

	read.creator = func() (interface{}, error) {
		if !found {
			return nil, fmt.Errorf("Unable to find an event of type %s", read.Type)
		}

		e := creator()
		err := json.Unmarshal(read.Content, &e)

		return e, err
	}

	return read, nil
}
