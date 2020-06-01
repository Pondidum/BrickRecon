package app

import (
	"fmt"
	"mvc/eventstore"
	"mvc/lego"
	"os"
)

type SiteModel struct {
	AllModels []string
}

type AppStore struct {
	es *eventstore.EventStore
}

func NewAppStore() (*AppStore, error) {
	if err := os.MkdirAll("_store", os.ModePerm); err != nil {
		return nil, err
	}

	es := eventstore.NewEventStore("_store")

	es.RegisterEvent(func() interface{} { return &lego.ProjectCreated{} })
	es.RegisterEvent(func() interface{} { return &lego.PartsAdded{} })

	es.RegisterProjection("projects", lego.ProjectsInitialState, lego.ProjectsProjector)

	return &AppStore{es: es}, nil
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
