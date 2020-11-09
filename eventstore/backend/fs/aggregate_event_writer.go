package fs

import (
	"brickrecon/eventstore"
	"bytes"
	"context"
	"encoding/json"
	"os"
	"path"

	"github.com/honeycombio/beeline-go"
)

type DirectoryPath string

type AggregateEventWriter struct {
	root DirectoryPath
}

func NewAggregateEventWriter(root DirectoryPath) *AggregateEventWriter {
	return &AggregateEventWriter{root}
}

func (ew *AggregateEventWriter) WriteEvents(ctx context.Context, aggregateID eventstore.AggregateID, changes []eventstore.Event) (int, error) {

	block := bytes.Buffer{}
	newline := []byte("\n")
	for _, e := range changes {

		bytes, err := json.Marshal(e)

		if err != nil {
			beeline.AddField(ctx, "es.event_serialization_err", err)
			return 0, err
		}

		block.Write(bytes)
		block.Write(newline)
	}

	filepath := path.Join(string(ew.root), string(aggregateID))

	beeline.AddField(ctx, "es.event_file", filepath)

	file, err := os.OpenFile(filepath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
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
