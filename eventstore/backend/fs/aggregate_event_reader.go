package fs

import (
	"brickrecon/eventstore"
	"bufio"
	"context"
	"encoding/json"
	"io"
	"os"
	"path"

	"github.com/honeycombio/beeline-go"
)

type AggregateEventReader struct {
	ctx      context.Context
	registry *eventstore.EventRegistry
	file     *os.File
	reader   *bufio.Reader

	line []byte
	err  error
}

func NewAggregateEventReader(ctx context.Context, registry *eventstore.EventRegistry, root DirectoryPath, id string) (*AggregateEventReader, error) {
	filepath := path.Join(string(root), id)

	file, err := os.Open(filepath)
	if err != nil && !os.IsNotExist(err) {
		beeline.AddField(ctx, "es.event_file_open_err", err)
		return nil, err
	}

	reader := bufio.NewReader(file)

	er := &AggregateEventReader{
		ctx:      ctx,
		registry: registry,
		file:     file,
		reader:   reader,
	}
	return er, nil
}

func (er *AggregateEventReader) Close() error {
	return er.file.Close()
}

func (er *AggregateEventReader) Read() bool {

	line, err := er.reader.ReadBytes('\n')

	if len(line) == 0 {
		return false
	}

	er.err = nil
	er.line = line

	if err == nil {
		return true
	}

	if err == io.EOF {
		return true
	}

	er.err = err

	return false
}

type aggregateEventType struct {
	Type string `json:"meta_type"`
}

func (er *AggregateEventReader) Event() (eventstore.Event, error) {

	if err := er.err; err != nil {
		return nil, err
	}

	var et aggregateEventType
	if err := json.Unmarshal(er.line, &et); err != nil {
		beeline.AddField(er.ctx, "es.event_type_unmarshal_err", err)
		return nil, err
	}

	event, err := er.registry.CreateInstance(er.ctx, et.Type)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(er.line, event); err != nil {
		beeline.AddField(er.ctx, "es.event_unmarshal_err", err)
		return nil, err
	}

	return event.(eventstore.Event), nil
}
