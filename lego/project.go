package lego

import uuid "github.com/satori/go.uuid"

type Project struct {
	id      uuid.UUID
	changes []interface{}
	version int

	Name  string
	parts *PartList
}

func NewProject(name string, parts []Part) *Project {

	project := Project{}
	project.apply(&ProjectCreated{ID: uuid.NewV4(), Name: name})
	project.apply(&PartsAdded{Parts: parts})

	return &project
}

func (prj *Project) FromEvents(events []interface{}) {
	for _, event := range events {
		prj.on(event)
		prj.version++
	}
}

func (prj *Project) apply(event interface{}) {
	prj.changes = append(prj.changes, event)
	prj.on(event)
}

func (prj *Project) ID() uuid.UUID {
	return prj.id
}

func (prj *Project) Version() int {
	return prj.version
}

func (prj *Project) Changes() []interface{} {
	return prj.changes
}

func (prj *Project) ClearChanges() {
	prj.changes = []interface{}{}
}

func (prj *Project) on(event interface{}) {

	switch e := event.(type) {

	case *ProjectCreated:
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
