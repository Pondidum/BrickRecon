package fs

import (
	"brickrecon/eventstore"
	"bytes"
	"context"
	"encoding/json"
	"os"

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

func (ew *FsEventWriter) WriteEvents(ctx context.Context, aggregateID uuid.UUID, changes []eventstore.Event) (int, error) {

	block := bytes.Buffer{}

	for _, e := range changes {

		bytes, err := json.Marshal(e)

		if err != nil {
			beeline.AddField(ctx, "es.event_serialization_err", err)
			return 0, err
		}

		block.Write(bytes)
		block.Write(newline)
	}

	beeline.AddField(ctx, "es.event_file", ew.filename)

	file, err := os.OpenFile(ew.filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
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

	return len(changes), nil
}
