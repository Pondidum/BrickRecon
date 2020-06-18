package app

import (
	"brickrecon/background"
	"brickrecon/distributor"
	"brickrecon/eventstore"
	"brickrecon/eventstore/backend/fs"
	"brickrecon/lego"
	"context"
	"os"
)

type AppStoreBuilder struct {
	ctx context.Context
}

func NewAppStoreBuilder(ctx context.Context) *AppStoreBuilder {
	return &AppStoreBuilder{ctx: ctx}
}

func (b *AppStoreBuilder) CreateBackend() (*fs.FsBackend, error) {
	if err := os.MkdirAll("_store", os.ModePerm); err != nil {
		return nil, err
	}

	backend, err := fs.NewFileSystemBackend("_store")
	if err != nil {
		return nil, err
	}

	return backend, err
}

func (b *AppStoreBuilder) CreateEventStore(backend eventstore.Backend) eventstore.EventStore {
	es := eventstore.NewEventStore(backend)

	lego.ProjectEvents(b.ctx, es.RegisterEvent)
	background.ImageCacheEvents(b.ctx, es.RegisterEvent)

	es.RegisterProjection(b.ctx, "projects", lego.ProjectsInitialState, lego.ProjectsProjector)

	return es
}

func (b *AppStoreBuilder) CreateBus(es eventstore.EventStore) *distributor.Distributor {

	bus := distributor.NewDistributor()

	bus.RegisterFor(&background.PartsAddedMessage{}, background.ImageCacheHandler(es))

	return bus
}

func (b *AppStoreBuilder) Create() (*AppStore, error) {

	backend, err := b.CreateBackend()
	if err != nil {
		return nil, err
	}

	es := b.CreateEventStore(backend)
	bus := b.CreateBus(es)

	return &AppStore{EventStore: es, bus: bus}, nil
}
