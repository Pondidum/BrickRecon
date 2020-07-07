package app

import (
	"brickrecon/distributor"
	"brickrecon/eventstore"
	"brickrecon/lego"
	"brickrecon/lego/projections/all_kits"
	"brickrecon/lego/projections/all_projects"
	"context"
	"fmt"
)

type SiteModel struct {
	AllKits   map[lego.KitNumber]*all_kits.KitView
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
	var projects all_projects.AllProjectsView
	if err := a.EventStore.ReadView(ctx, all_projects.ProjectionName, &projects); err != nil {
		return SiteModel{}
	}

	var kits all_kits.AllKitsView
	if err := a.EventStore.ReadView(ctx, all_kits.ProjectionName, &kits); err != nil {
		return SiteModel{}
	}

	return SiteModel{
		AllModels: projects.Names,
		AllKits:   kits.Kits,
	}
}

func (a *AppStore) ReadProject(ctx context.Context, name lego.ProjectName) (*all_projects.ProjectView, error) {

	var view all_projects.AllProjectsView
	if err := a.EventStore.ReadView(ctx, all_projects.ProjectionName, &view); err != nil {
		return nil, err
	}

	project, ok := view.Projects[name]

	if !ok {
		return nil, fmt.Errorf("No project with the name '%s' was found", name)
	}

	return project, nil

}

func (a *AppStore) ReadKit(ctx context.Context, kitNumber lego.KitNumber) (*all_kits.KitView, error) {

	var view all_kits.AllKitsView
	if err := a.EventStore.ReadView(ctx, all_kits.ProjectionName, &view); err != nil {
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
