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

func NewFileSystemBackend(root string) (*FsBackend, error) {

	be := &FsBackend{root: root}

	if err := be.createRoot(); err != nil {
		return nil, err
	}

	return be, nil
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

func (be *FsBackend) DestroyViews() error {

	if err := os.RemoveAll(path.Join(be.root, "views")); err != nil {
		return err
	}

	return be.createRoot()
}

func (be *FsBackend) createRoot() error {
	return os.MkdirAll(path.Join(be.root, "views"), os.ModePerm)
}
