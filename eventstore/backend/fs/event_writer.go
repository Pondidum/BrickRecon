package fs

import (
	"brickrecon/eventstore"
	"bytes"
	"context"
	"encoding/json"
	"os"
	"time"

	"github.com/honeycombio/beeline-go"
	uuid "github.com/satori/go.uuid"
)

var newline = []byte("\n")

type FsEventWriter struct {
	filename string
}

func NewEventWriter(filename string) *FsEventWriter {
	return &FsEventWriter{filename}
}

func (ew *FsEventWriter) WriteEvents(ctx context.Context, aggregateID uuid.UUID, currentVersion int, changes []eventstore.Event) (int, error) {

	block := bytes.Buffer{}

	for _, e := range changes {

		currentVersion++

		meta := e.Meta()

		meta.Timestamp = time.Now()
		meta.ID = uuid.NewV4()
		meta.AggregateRootID = aggregateID
		meta.Version = currentVersion
		meta.Type = eventstore.EventName(e)

		bytes, err := json.Marshal(e)

		if err != nil {
			beeline.AddField(ctx, "es.event_serialization_err", err)
			return 0, err
		}

		block.Write(bytes)
		block.Write(newline)
	}

	file, err := os.OpenFile(ew.filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	beeline.AddField(ctx, "es.event_file", ew.filename)

	if err != nil {
		beeline.AddField(ctx, "es.event_file_open_err", err)
		return 0, err
	}
	defer file.Close()

	if _, err := file.Write(block.Bytes()); err != nil {
		beeline.AddField(ctx, "es.event_file_write_err", err)
		return 0, err
	}

	beeline.AddField(ctx, "es.events_written_count", len(changes))

	return currentVersion, nil
}
