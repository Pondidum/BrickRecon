package fs

import (
	"brickrecon/eventstore"
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path"

	"github.com/honeycombio/beeline-go"
	uuid "github.com/satori/go.uuid"
)

type AggregateEventReader struct {
	ctx      context.Context
	registry map[string]eventstore.Initialiser
	file     *os.File
	scanner  *bufio.Scanner
}

func NewAggregateEventReader(ctx context.Context, registry map[string]eventstore.Initialiser, root DirectoryPath, id uuid.UUID) (*AggregateEventReader, error) {
	filepath := path.Join(string(root), id.String())

	file, err := os.Open(filepath)
	if err != nil && !os.IsNotExist(err) {
		beeline.AddField(ctx, "es.event_file_open_err", err)
		return nil, err
	}

	scanner := bufio.NewScanner(file)

	return &AggregateEventReader{ctx, registry, file, scanner}, nil
}

func (er *AggregateEventReader) Close() error {
	return er.file.Close()
}

func (er *AggregateEventReader) Read() bool {
	return er.scanner.Scan()
}

type aggregateEventType struct {
	Type string `json:"meta_type"`
}

func (er *AggregateEventReader) Event() (eventstore.Event, error) {
	var et aggregateEventType
	if err := json.Unmarshal(er.scanner.Bytes(), &et); err != nil {
		beeline.AddField(er.ctx, "es.event_type_unmarshal_err", err)
		return nil, err
	}

	creator, found := er.registry[et.Type]

	if !found {
		beeline.AddField(er.ctx, "es.event_type_lookup_failed", true)
		return nil, fmt.Errorf("Unable to find an event of type %s", et.Type)
	}

	event := creator()
	if err := json.Unmarshal(er.scanner.Bytes(), event); err != nil {
		beeline.AddField(er.ctx, "es.event_unmarshal_err", err)
		return nil, err
	}

	return event.(eventstore.Event), nil
}
