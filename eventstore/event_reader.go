package eventstore

import (
	"bufio"
	"encoding/json"
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
	Content     json.RawMessage
}

func (r *Event) Event(e interface{}) error {
	return json.Unmarshal(r.Content, &e)
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

func (er *EventReader) Record() (Event, error) {
	var read Event
	err := json.Unmarshal(er.scanner.Bytes(), &read)
	return read, err
}
