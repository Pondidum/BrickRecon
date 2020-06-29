package app

import (
	"brickrecon/distributor"
	"brickrecon/eventstore"
	"brickrecon/lego"
	"context"
	"fmt"
)

type SiteModel struct {
	AllKits   map[lego.KitNumber]*lego.KitView
	AllModels []lego.ProjectName
}

type AppStore struct {
	EventStore eventstore.EventStore
	bus        *distributor.Distributor
}

func (a *AppStore) Save(ctx context.Context, aggregate eventstore.Aggregate) error {
	return a.EventStore.SaveAggregate(ctx, aggregate)
}

func (a *AppStore) SiteModel(ctx context.Context) SiteModel {
	var projects lego.AllProjectsView
	if err := a.EventStore.ReadView(ctx, "projects", &projects); err != nil {
		return SiteModel{}
	}

	var kits lego.AllKitsView
	if err := a.EventStore.ReadView(ctx, "kits", &kits); err != nil {
		return SiteModel{}
	}

	return SiteModel{
		AllModels: projects.Names,
		AllKits:   kits.Kits,
	}
}

func (a *AppStore) ReadProject(ctx context.Context, name lego.ProjectName) (*lego.ProjectView, error) {

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

func (a *AppStore) ReadKit(ctx context.Context, kitNumber lego.KitNumber) (*lego.KitView, error) {

	var view lego.AllKitsView
	if err := a.EventStore.ReadView(ctx, "kits", &view); err != nil {
		return nil, err
	}

	kit, ok := view.Kits[kitNumber]

	if !ok {
		return nil, fmt.Errorf("No kit with kitnumber '%s' was found", kitNumber)
	}

	return kit, nil
}

func (a *AppStore) SendMessage(ctx context.Context, message distributor.Message) func() {
	return a.bus.Dispatch(ctx, message)
}
