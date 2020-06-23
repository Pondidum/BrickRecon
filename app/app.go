package app

import (
	"brickrecon/distributor"
	"brickrecon/eventstore"
	"brickrecon/lego"
	"context"
	"fmt"
)

type SiteModel struct {
	AllModels []string
}

type AppStore struct {
	EventStore eventstore.EventStore
	bus        *distributor.Distributor
}

func (a *AppStore) Save(ctx context.Context, aggregate eventstore.Aggregate) error {
	return a.EventStore.SaveAggregate(ctx, aggregate)
}

func (a *AppStore) SiteModel(ctx context.Context) SiteModel {
	var view lego.AllProjectsView
	if err := a.EventStore.ReadView(ctx, "projects", &view); err != nil {
		return SiteModel{}
	}

	return SiteModel{AllModels: view.Names}
}

func (a *AppStore) ReadProject(ctx context.Context, name string) (*lego.ProjectView, error) {

	var view lego.AllProjectsView
	if err := a.EventStore.ReadView(ctx, "projects", &view); err != nil {
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
