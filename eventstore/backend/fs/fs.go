package fs

import (
	"brickrecon/eventstore"
	"context"
	"os"
	"path"
)

type FsBackend struct {
	root string
}

func NewFileSystemBackend(root string) (eventstore.Backend, error) {

	if err := os.MkdirAll(path.Join(root, "views"), os.ModePerm); err != nil {
		return nil, err
	}

	return &FsBackend{
		root: root,
	}, nil
}

func (be *FsBackend) NewEventReader(registry map[string]eventstore.Initialiser, ctx context.Context) (eventstore.EventReader, error) {
	return NewEventReader(registry, path.Join(be.root, "events"), ctx)
}

func (be *FsBackend) NewEventWriter() eventstore.EventWriter {
	return NewEventWriter(path.Join(be.root, "events"))
}

func (be *FsBackend) NewView(name string) eventstore.View {
	return &FsView{
		filename: path.Join(be.root, "views", name+".json"),
	}
}
