package eventstore

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
)

type EventReader struct {
	registry map[string]Initialiser
	file     *os.File
	scanner  *bufio.Scanner
}

func NewEventReader(registry map[string]Initialiser, filename string) (*EventReader, error) {
	file, err := os.Open(filename)

	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(file)

	return &EventReader{registry, file, scanner}, nil
}

func (er *EventReader) Close() error {
	return er.file.Close()
}

func (er *EventReader) ReadAll() bool {
	return er.scanner.Scan()
}

func (er *EventReader) Event() (interface{}, error) {
	var read readEvent
	if err := json.Unmarshal(er.scanner.Bytes(), &read); err != nil {
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
