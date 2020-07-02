package fs

import (
	"brickrecon/eventstore"
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"

	uuid "github.com/satori/go.uuid"
)

type FsEventReader struct {
	registry     map[string]eventstore.Initialiser
	file         *os.File
	scanner      *bufio.Scanner
	ctx          context.Context
	currentIndex int
}

func NewEventReader(registry map[string]eventstore.Initialiser, filename string, ctx context.Context) (*FsEventReader, error) {
	file, err := os.Open(filename)

	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	scanner := bufio.NewScanner(file)

	return &FsEventReader{registry, file, scanner, ctx, 0}, nil
}

func (er *FsEventReader) Close() error {
	return er.file.Close()
}

func (er *FsEventReader) ReadAll() bool {
	er.currentIndex++
	return er.scanner.Scan()
}

func (er *FsEventReader) ReadFor(uuid uuid.UUID) bool {

	for er.scanner.Scan() {
		er.currentIndex++
		record, err := er.Event()

		if err != nil {
			return false
		}

		if record.Meta().AggregateRootID == uuid {
			return true
		}
	}

	return false
}

func (er *FsEventReader) ReadFrom(offset int) bool {

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

func (er *FsEventReader) Event() (eventstore.Event, error) {
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

	return event.(eventstore.Event), nil
}
