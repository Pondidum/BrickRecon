package app

import (
	"brickrecon/background"
	"brickrecon/distributor"
	"brickrecon/eventstore"
	"brickrecon/eventstore/backend/fs"
	"brickrecon/lego"
	"brickrecon/preen"
	"context"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/honeycombio/beeline-go/wrappers/hnynethttp"
)

type AppBuilder struct {
	ctx context.Context
}

func NewAppBuilder(ctx context.Context) *AppBuilder {
	return &AppBuilder{ctx: ctx}
}

func (b *AppBuilder) CreateBackend() (*fs.FsBackend, error) {
	if err := os.MkdirAll("_store", os.ModePerm); err != nil {
		return nil, err
	}

	backend, err := fs.NewFileSystemBackend("_store")
	if err != nil {
		return nil, err
	}

	return backend, err
}

func (b *AppBuilder) CreateEventStore(backend eventstore.Backend) eventstore.EventStore {
	es := eventstore.NewEventStore(backend)

	eventstore.RegisterMany(es, b.ctx, lego.ProjectEvents)
	eventstore.RegisterMany(es, b.ctx, lego.KitEvents)
	eventstore.RegisterMany(es, b.ctx, background.ImageCacheEvents)

	es.RegisterProjection(b.ctx, &lego.ProjectsProjection{})
	es.RegisterProjection(b.ctx, &lego.KitsProjection{})
	es.RegisterProjection(b.ctx, &lego.ProjectKitsProjection{})

	return es
}

func (b *AppBuilder) CreateBus(es eventstore.EventStore) *distributor.Distributor {

	bus := distributor.NewDistributor()

	bus.RegisterFor(&background.PartsAddedMessage{}, background.ImageCacheHandler(es))

	return bus
}

func (b *AppBuilder) CreateAppStore() (*AppStore, error) {

	backend, err := b.CreateBackend()
	if err != nil {
		return nil, err
	}

	es := b.CreateEventStore(backend)
	bus := b.CreateBus(es)

	return &AppStore{EventStore: es, bus: bus}, nil
}

func (b *AppBuilder) CreateWebUI() (http.Handler, error) {

	store, err := b.CreateAppStore()

	if err != nil {
		return nil, err
	}

	p, err := preen.NewPreen(preen.PreenConfig{
		ApplicationRoot: "app",
		Controllers: []preen.Controller{
			&RootController{Store: store},
			&CreateController{Store: store},
			&ProjectController{Store: store},
			&LoginController{},
			&KitImportController{Store: store},
			&KitController{Store: store},
		},
	})

	if err != nil {
		return nil, err
	}

	r := mux.NewRouter()

	p.Apply(r)

	return hnynethttp.WrapHandler(r), nil
}
