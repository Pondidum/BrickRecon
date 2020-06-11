package eventstore

import (
	"bytes"
	"encoding/json"
	"os"
	"time"

	uuid "github.com/satori/go.uuid"
)

type EventWriter interface {
	WriteEvents(aggregateID uuid.UUID, currentVersion int, changes []Event) (int, error)
}

type FsEventWriter struct {
	filename string
}

func NewEventWriter(filename string) EventWriter {
	return &FsEventWriter{filename}
}

func (ew *FsEventWriter) WriteEvents(aggregateID uuid.UUID, currentVersion int, changes []Event) (int, error) {
	file, err := os.OpenFile(ew.filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)

	if err != nil {
		return 0, err
	}

	defer file.Close()

	block := bytes.Buffer{}

	for _, e := range changes {

		currentVersion++

		meta := e.event()

		meta.Timestamp = time.Now()
		meta.ID = uuid.NewV4()
		meta.AggregateRootID = aggregateID
		meta.Version = currentVersion
		meta.Type = eventName(e)

		bytes, err := json.Marshal(e)

		if err != nil {
			return 0, err
		}

		if _, err := block.Write(bytes); err != nil {
			return 0, err
		}

		if _, err := block.Write(newline); err != nil {
			return 0, err
		}
	}

	if _, err := file.Write(block.Bytes()); err != nil {
		return 0, err
	}

	return currentVersion, nil
}
