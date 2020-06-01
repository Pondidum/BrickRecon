package lego

import (
	"mvc/eventstore"

	uuid "github.com/satori/go.uuid"
)

type AllProjectsView struct {
	Names    []string
	Projects map[string]*ProjectView
}

type ProjectView struct {
	ID   uuid.UUID
	Name string

	Parts []Part
}

func ProjectsInitialState() interface{} {
	return &AllProjectsView{
		Names:    []string{},
		Projects: map[string]*ProjectView{},
	}
}

func ProjectsProjector(state interface{}, record eventstore.Record) interface{} {
	view := state.(*AllProjectsView)

	switch e := record.Event().(type) {
	case ProjectCreated:
		view.Names = append(view.Names, e.Name)
		view.Projects[e.Name] = &ProjectView{ID: record.AggregateID, Name: e.Name}

	case PartsAdded:
		project := projectByID(view.Projects, record.AggregateID)
		project.Parts = append(project.Parts, e.Parts...)
	}

	return view
}

func projectByID(all map[string]*ProjectView, id uuid.UUID) *ProjectView {
	for _, p := range all {
		if p.ID == id {
			return p
		}
	}
	return nil
}
