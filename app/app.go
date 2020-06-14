package app

import (
	"brickrecon/background"
	"brickrecon/distributor"
	"brickrecon/eventstore"
	"brickrecon/eventstore/backend/fs"
	"brickrecon/lego"
	"context"
	"fmt"
	"os"
)

type SiteModel struct {
	AllModels []string
}

type AppStore struct {
	es  eventstore.EventStore
	bus *distributor.Distributor
}

func NewAppStore() (*AppStore, error) {
	if err := os.MkdirAll("_store", os.ModePerm); err != nil {
		return nil, err
	}

	backend, err := fs.NewFileSystemBackend("_store")
	if err != nil {
		return nil, err
	}

	es := eventstore.NewEventStore(backend)

	lego.ProjectEvents(es.RegisterEvent)
	background.ImageCacheEvents(es.RegisterEvent)

	es.RegisterProjection("projects", lego.ProjectsInitialState, lego.ProjectsProjector)

	bus := distributor.NewDistributor()

	if err := background.AttachImageCacheListener(bus, es); err != nil {
		return nil, err
	}

	return &AppStore{es: es, bus: bus}, nil
}

func (a *AppStore) Save(project *lego.Project) error {
	return a.es.SaveAggregate(project.Aggregator)
}

func (a *AppStore) SiteModel() SiteModel {
	var view lego.AllProjectsView
	if err := a.es.ReadView("projects", &view); err != nil {
		return SiteModel{}
	}

	return SiteModel{AllModels: view.Names}
}

func (a *AppStore) Project(name string) (*lego.ProjectView, error) {

	var view lego.AllProjectsView
	if err := a.es.ReadView("projects", &view); err != nil {
		return nil, err
	}

	project, ok := view.Projects[name]

	if !ok {
		return nil, fmt.Errorf("No project with the name '%s' was found", name)
	}

	return project, nil

}

func (a *AppStore) SendMessage(ctx context.Context, message distributor.Message) func() {
	return a.bus.Dispatch(ctx, message)
}
