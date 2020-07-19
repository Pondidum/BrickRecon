package fs

import (
	"brickrecon/eventstore"
	"context"
	"io/ioutil"
	"os"
	"path"

	uuid "github.com/satori/go.uuid"
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

func (be *AggregateBackend) NewEventReader(registry map[string]eventstore.Initialiser, aggregateID uuid.UUID) (eventstore.EventReader, error) {
	return NewAggregateEventReader(context.Background(), registry, be.eventsPath, aggregateID)
}

func (be *AggregateBackend) NewEventWriter() eventstore.EventWriter {
	return NewAggregateEventWriter(be.eventsPath)
}

func (be *AggregateBackend) NewView(name string) eventstore.View {
	return &FsView{
		filename: path.Join(string(be.viewsPath), name+".json"),
	}
}

func (be *AggregateBackend) DestroyViews() error {

	if err := os.RemoveAll(string(be.viewsPath)); err != nil {
		return err
	}

	return be.createRoot()
}

func (be *AggregateBackend) AllAggregates() ([]uuid.UUID, error) {

	entries, err := ioutil.ReadDir(string(be.eventsPath))
	if err != nil {
		return nil, err
	}

	ids := make([]uuid.UUID, len(entries))

	for i, info := range entries {
		id, err := uuid.FromString(info.Name())
		if err != nil {
			return nil, err
		}

		ids[i] = id
	}

	return ids, nil
}

// type aggregateID struct {
// 	AggregateRootID uuid.UUID `json:"meta_aggregate_id"`
// }

func (be *AggregateBackend) createRoot() error {
	if err := os.MkdirAll(string(be.eventsPath), os.ModePerm); err != nil {
		return err
	}

	if err := os.MkdirAll(string(be.viewsPath), os.ModePerm); err != nil {
		return err
	}

	return nil
}
