package app

import (
	"brickrecon/distributor"
	"brickrecon/eventstore"
	"brickrecon/lego"
	"brickrecon/lego/projections/allkits"
	"brickrecon/lego/projections/allprojects"
	"context"
	"fmt"
)

type SiteModel struct {
	AllKits   map[lego.KitNumber]*allkits.KitView
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
	var projects allprojects.AllProjectsView
	if err := a.EventStore.ReadView(ctx, allprojects.ProjectionName, &projects); err != nil {
		return SiteModel{}
	}

	var kits allkits.AllKitsView
	if err := a.EventStore.ReadView(ctx, allkits.ProjectionName, &kits); err != nil {
		return SiteModel{}
	}

	return SiteModel{
		AllModels: projects.Names,
		AllKits:   kits.Kits,
	}
}

func (a *AppStore) ReadProjectView(ctx context.Context, name lego.ProjectName) (*allprojects.ProjectView, error) {

	var view allprojects.AllProjectsView
	if err := a.EventStore.ReadView(ctx, allprojects.ProjectionName, &view); err != nil {
		return nil, err
	}

	project, ok := view.Projects[name]

	if !ok {
		return nil, fmt.Errorf("No project with the name '%s' was found", name)
	}

	return project, nil

}

func (a *AppStore) ReadProject(ctx context.Context, name lego.ProjectName) (*lego.Project, error) {

	selected, err := a.ReadProjectView(ctx, name)
	if err != nil {
		return nil, err
	}

	project := lego.BlankProject()
	if err := a.EventStore.LoadAggregate(ctx, selected.ID, project); err != nil {
		return nil, err
	}

	return project, nil
}

func (a *AppStore) ReadProjectByID(ctx context.Context, id eventstore.AggregateID) (*lego.Project, error) {

	project := lego.BlankProject()
	if err := a.EventStore.LoadAggregate(ctx, id, project); err != nil {
		return nil, err
	}

	return project, nil
}

func (a *AppStore) ReadKitView(ctx context.Context, kitNumber lego.KitNumber) (*allkits.KitView, error) {

	var view allkits.AllKitsView
	if err := a.EventStore.ReadView(ctx, allkits.ProjectionName, &view); err != nil {
		return nil, err
	}

	kit, ok := view.Kits[kitNumber]

	if !ok {
		return nil, fmt.Errorf("No kit with kitnumber '%s' was found", kitNumber)
	}

	return kit, nil
}

func (a *AppStore) ReadPart(ctx context.Context, key lego.PartKey) (*lego.PartA, error) {

	part := lego.BlankPart()
	if err := a.EventStore.LoadAggregate(ctx, eventstore.AggregateID(key), part); err != nil {
		return nil, err
	}

	return part, nil
}

func (a *AppStore) SendMessage(ctx context.Context, message distributor.Message) func() {
	return a.bus.Dispatch(ctx, message)
}
