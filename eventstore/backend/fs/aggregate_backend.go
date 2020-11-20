package fs

import (
	"brickrecon/eventstore"
	"context"
	"io/ioutil"
	"os"
	"path"

	"github.com/honeycombio/beeline-go"
)

type AggregateBackend struct {
	eventsPath DirectoryPath
	viewsPath  DirectoryPath
}

func NewAggregateBackend(root string) (*AggregateBackend, error) {

	be := &AggregateBackend{
		eventsPath: DirectoryPath(path.Join(root, "events")),
		viewsPath:  DirectoryPath(path.Join(root, "views")),
	}

	if err := be.createRoot(); err != nil {
		return nil, err
	}

	return be, nil
}

func (be *AggregateBackend) NewEventReader(ctx context.Context, registry *eventstore.EventRegistry, aggregateID eventstore.AggregateID) (eventstore.EventReader, error) {
	return NewAggregateEventReader(ctx, registry, be.eventsPath, string(aggregateID))
}

func (be *AggregateBackend) NewEventWriter() eventstore.EventWriter {
	return NewAggregateEventWriter(be.eventsPath)
}

func (be *AggregateBackend) NewView(name string) eventstore.View {
	return &FsView{
		filename: path.Join(string(be.viewsPath), name+".json"),
	}
}

func (be *AggregateBackend) DestroyViews(ctx context.Context) error {
	ctx, span := beeline.StartSpan(ctx, "destroy_views")
	defer span.Send()

	if err := os.RemoveAll(string(be.viewsPath)); err != nil {
		return err
	}

	return be.createRoot()
}

func (be *AggregateBackend) AllAggregates(ctx context.Context) ([]eventstore.AggregateID, error) {

	ctx, span := beeline.StartSpan(ctx, "all_aggregates")
	defer span.Send()

	entries, err := ioutil.ReadDir(string(be.eventsPath))
	if err != nil {
		return nil, err
	}

	ids := make([]eventstore.AggregateID, len(entries))

	for i, info := range entries {
		ids[i] = eventstore.AggregateID(info.Name())
	}

	return ids, nil
}

func (be *AggregateBackend) createRoot() error {
	if err := os.MkdirAll(string(be.eventsPath), os.ModePerm); err != nil {
		return err
	}

	if err := os.MkdirAll(string(be.viewsPath), os.ModePerm); err != nil {
		return err
	}

	return nil
}
