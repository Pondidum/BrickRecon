package lego

import (
	"mvc/eventstore"
	"reflect"

	"github.com/honeycombio/libhoney-go"
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

func ProjectsProjector(state interface{}, event eventstore.Event) interface{} {
	view := state.(*AllProjectsView)

	o := libhoney.NewEvent()
	defer o.Send()

	o.AddField("event_type", reflect.TypeOf(event).Elem().Name())

	switch e := event.(type) {
	case *ProjectCreated:
		o.AddField("project_name", e.Name)
		o.AddField("project_id", e.AggregateID())

		view.Names = append(view.Names, e.Name)
		view.Projects[e.Name] = &ProjectView{ID: e.AggregateID(), Name: e.Name}

	case *ProjectPartsAdded:
		o.AddField("project_id", e.AggregateRootID)
		o.AddField("part_count", len(e.Parts))

		project := projectByID(view.Projects, e.AggregateID())
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
