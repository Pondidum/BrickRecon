package fs

import (
	"brickrecon/eventstore"
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"

	"github.com/honeycombio/beeline-go"
	uuid "github.com/satori/go.uuid"
)

type AggregateEventReader struct {
	ctx      context.Context
	registry map[string]eventstore.Initialiser
	file     *os.File
	reader   *bufio.Reader

	line []byte
	err  error
}

func NewAggregateEventReader(ctx context.Context, registry map[string]eventstore.Initialiser, root DirectoryPath, id uuid.UUID) (*AggregateEventReader, error) {
	filepath := path.Join(string(root), id.String())

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

	creator, found := er.registry[et.Type]

	if !found {
		beeline.AddField(er.ctx, "es.event_type_lookup_failed", true)
		return nil, fmt.Errorf("Unable to find an event of type %s", et.Type)
	}

	event := creator()
	if err := json.Unmarshal(er.line, event); err != nil {
		beeline.AddField(er.ctx, "es.event_unmarshal_err", err)
		return nil, err
	}

	return event.(eventstore.Event), nil
}
