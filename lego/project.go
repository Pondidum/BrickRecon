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

	project := Project{
		parts: &PartList{},
	}
	project.Aggregator = eventstore.NewAggregator(project.on)

	project.Apply(&ProjectCreated{ID: uuid.NewV4(), Name: name})
	project.Apply(&ProjectPartsAdded{Parts: parts})

	return &project
}

func (prj *Project) on(event eventstore.Event) {

	switch e := event.(type) {

	case *ProjectCreated:
		prj.SetID(e.ID)
		prj.Name = e.Name

	case *ProjectPartsAdded:
		for _, p := range e.Parts {
			prj.parts.Add(p)
		}
	}

}

type ProjectCreated struct {
	eventstore.EventMeta

	ID   uuid.UUID
	Name string
}

type ProjectPartsAdded struct {
	eventstore.EventMeta

	Parts []Part
}

func ProjectEvents(register func(eventstore.Initialiser)) {
	register(func() interface{} { return &ProjectCreated{} })
	register(func() interface{} { return &ProjectPartsAdded{} })
}
