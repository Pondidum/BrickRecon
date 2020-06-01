package lego

import (
	"mvc/eventstore"

	uuid "github.com/satori/go.uuid"
)

type Project struct {
	*eventstore.Aggregator

	Name  string
	parts *PartList
}

func NewProject(name string, parts []Part) *Project {

	project := Project{}
	project.Aggregator = eventstore.NewAggregator(project.on)

	project.Apply(&ProjectCreated{ID: uuid.NewV4(), Name: name})
	project.Apply(&PartsAdded{Parts: parts})

	return &project
}

func (prj *Project) on(event interface{}) {

	switch e := event.(type) {

	case *ProjectCreated:
		prj.SetID(e.ID)
		prj.Name = e.Name

	case *PartsAdded:
		for _, p := range e.Parts {
			prj.parts.Add(p)
		}
	}

}

type ProjectCreated struct {
	ID   uuid.UUID
	Name string
}

type PartsAdded struct {
	Parts []Part
}
