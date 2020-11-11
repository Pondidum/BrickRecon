package app

import (
	"brickrecon/background"
	"brickrecon/distributor"
	"brickrecon/eventstore"
	"brickrecon/eventstore/backend/fs"
	"brickrecon/lego"
	"brickrecon/lego/projections/all_kits"
	"brickrecon/lego/projections/all_projects"
	"brickrecon/lego/projections/allparts"
	"brickrecon/lego/projections/colours"
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

func (b *AppBuilder) CreateBackend() (*fs.AggregateBackend, error) {
	if err := os.MkdirAll("_store", os.ModePerm); err != nil {
		return nil, err
	}

	backend, err := fs.NewAggregateBackend("_store")
	if err != nil {
		return nil, err
	}

	return backend, err
}

func (b *AppBuilder) CreateEventStore(backend eventstore.Backend) eventstore.EventStore {
	es := eventstore.NewEventStore(backend)

	es.RegisterEvents(b.ctx, lego.ProjectEvents)
	es.RegisterEvents(b.ctx, lego.KitEvents)
	es.RegisterEvents(b.ctx, lego.PartEvents)
	es.RegisterEvents(b.ctx, background.ImageCacheEvents)

	es.RegisterProjection(b.ctx, &all_projects.ProjectsProjection{})
	es.RegisterProjection(b.ctx, &all_kits.KitsProjection{})
	es.RegisterProjection(b.ctx, &colours.ColoursProjection{})
	es.RegisterProjection(b.ctx, &allparts.AllPartsProjection{})

	es.RegisterEventMiddleware(b.ctx, b.upgradeEvent)

	return es
}

func (b *AppBuilder) upgradeEvent(ctx context.Context, e eventstore.Event) eventstore.Event {

	switch event := e.(type) {

	case *lego.ProjectInventoryRemoved:
		if event.EventVersion == 0 {

			event.Part = lego.CreatePartKey(event.PartID, lego.LookupLDrawColour(int(event.ColourID)).LDrawID)
			event.EventVersion = 1
		}

	case *lego.ProjectInventoryAdded:
		if event.EventVersion == 0 {
			event.Part = lego.CreatePartKey(event.PartID, lego.LookupLDrawColour(int(event.ColourID)).LDrawID)
			event.EventVersion = 1
		}

	case *lego.ProjectPartsAdded:
		if event.EventVersion == 0 {
			for _, part := range event.Parts {
				part.Key = lego.CreatePartKey(part.Aliases.LDrawID, part.Colour.Aliases.LDrawID)
			}
			event.EventVersion = 1
		}

	case *lego.PartsChanged:
		if event.EventVersion == 0 {
			for _, part := range event.Additions {
				part.Key = lego.CreatePartKey(part.Aliases.LDrawID, part.Colour.Aliases.LDrawID)
			}

			replacement := map[lego.PartKey]int{}

			for oldKey, quantity := range event.Removals {
				id, oldColour := lego.ParsePartKey(oldKey)
				newKey := lego.CreatePartKey(id, lego.LookupLDrawColour(int(oldColour)).LDrawID)
				replacement[newKey] = quantity
			}

			event.Removals = replacement
			event.EventVersion = 1
		}
	}

	return e
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
		GetSiteModel: func(ctx context.Context) interface{} {
			return store.SiteModel(ctx)
		},
		Controllers: []preen.Controller{
			&RootController{Store: store},
			&CreateController{Store: store},
			&ProjectController{Store: store},
			&ProjectExportController{Store: store},
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
