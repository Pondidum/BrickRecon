package lego

import (
	"brickrecon/eventstore"

	uuid "github.com/satori/go.uuid"
)

type AllProjectsView struct {
	Names    []string
	Projects map[string]*ProjectView
}

type ProjectView struct {
	ID   uuid.UUID
	Name string

	Parts []*ProjectPartView
}

type ProjectPartView struct {
	ID         PartID
	Name       PartName
	ColourID   BrickLinkColour
	ColourName ColourName

	Quantity  int
	Inventory int
}

func toProjectPartView(part Part) *ProjectPartView {
	return &ProjectPartView{
		ID:         part.ID,
		Name:       part.Name,
		ColourID:   part.Colour.ID,
		ColourName: part.Colour.Name,
		Quantity:   part.Quantity,
	}
}

func ProjectsInitialState() interface{} {
	return &AllProjectsView{
		Names:    []string{},
		Projects: map[string]*ProjectView{},
	}
}

func ProjectsProjector(state interface{}, event eventstore.Event) interface{} {
	view := state.(*AllProjectsView)

	switch e := event.(type) {

	case *ProjectCreated:
		view.Names = append(view.Names, e.Name)
		view.Projects[e.Name] = &ProjectView{ID: e.AggregateID(), Name: e.Name}

	case *ProjectPartsAdded:
		project := projectByID(view.Projects, e.AggregateID())
		for _, part := range e.Parts {
			project.Parts = append(project.Parts, toProjectPartView(part))
		}

	case *ProjectInventoryAdded:
		project := projectByID(view.Projects, e.AggregateID())
		part := findPart(project.Parts, e.PartID, e.ColourID)
		part.Inventory += e.Quantity

	case *ProjectInventoryRemoved:
		project := projectByID(view.Projects, e.AggregateID())
		part := findPart(project.Parts, e.PartID, e.ColourID)
		part.Inventory -= e.Quantity

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

func findPart(parts []*ProjectPartView, partID PartID, colourID BrickLinkColour) *ProjectPartView {

	for _, part := range parts {
		if part.ID == partID && part.ColourID == colourID {
			return part
		}
	}

	return nil
}
