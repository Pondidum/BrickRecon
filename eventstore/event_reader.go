package eventstore

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"

	uuid "github.com/satori/go.uuid"
)

type EventReader struct {
	registry     map[string]Initialiser
	file         *os.File
	scanner      *bufio.Scanner
	currentIndex int
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
		read, err := er.readEvent()

		if err != nil {
			return false
		}

		if read.AggregateID == uuid {
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

func (er *EventReader) readEvent() (*readEvent, error) {
	var read readEvent
	if err := json.Unmarshal(er.scanner.Bytes(), &read); err != nil {
		return nil, err
	}

	return &read, nil
}

func (er *EventReader) Event() (interface{}, error) {

	read, err := er.readEvent()
	if err != nil {
		return nil, err
	}

	creator, found := er.registry[read.Type]

	if !found {
		return nil, fmt.Errorf("Unable to find an event of type %s", read.Type)
	}

	e := creator()

	if err := json.Unmarshal(read.Content, &e); err != nil {
		return nil, err
	}

	return e, nil
}
