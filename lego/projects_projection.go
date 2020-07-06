package lego

import (
	"brickrecon/eventstore"
	"fmt"

	uuid "github.com/satori/go.uuid"
)

type PartKey string

func CreatePartKey(part LDrawPart, colour BrickLinkColour) PartKey {
	return PartKey(fmt.Sprintf("%v|%v", part, colour))
}

type AllProjectsView struct {
	Names    []ProjectName
	Projects map[ProjectName]*ProjectView
	Kits     map[KitNumber]map[PartKey]int
}

type ProjectView struct {
	ID   uuid.UUID
	Name ProjectName

	Parts []*ProjectPartView

	Kits map[KitNumber]map[PartKey]int
}

type ProjectPartView struct {
	ID         LDrawPart
	Name       PartName
	ColourID   BrickLinkColour
	ColourName ColourName
	ColourHex  HexColour

	Key PartKey

	Quantity  int
	Inventory int
}

func toProjectPartView(part Part) *ProjectPartView {
	return &ProjectPartView{
		ID:         part.ID,
		Name:       part.Name,
		ColourID:   part.Colour.ID,
		ColourName: part.Colour.Name,
		ColourHex:  part.Colour.Hex,
		Quantity:   part.Quantity,
		Key:        CreatePartKey(part.ID, part.Colour.ID),
	}
}

var ProjectsProjectionName string = "projects"

type ProjectsProjection struct{}

func (p *ProjectsProjection) Name() string {
	return ProjectsProjectionName
}

func (p *ProjectsProjection) CreateState() interface{} {
	return &AllProjectsView{
		Names:    []ProjectName{},
		Projects: map[ProjectName]*ProjectView{},
		Kits:     map[KitNumber]map[PartKey]int{},
	}
}

func (p *ProjectsProjection) Project(state interface{}, event eventstore.Event) interface{} {
	view := state.(*AllProjectsView)

	switch e := event.(type) {

	case *ProjectCreated:
		view.Names = append(view.Names, e.Name)
		view.Projects[e.Name] = &ProjectView{
			ID:   e.AggregateRootID,
			Name: e.Name,
			Kits: map[KitNumber]map[PartKey]int{},
		}

	case *ProjectPartsAdded:
		project := projectByID(view.Projects, e.AggregateRootID)
		for _, part := range e.Parts {
			project.Parts = append(project.Parts, toProjectPartView(part))
		}

		for kn, kit := range view.Kits {
			calculateKitFulfillment(project, kn, kit)
		}

	case *ProjectInventoryAdded:
		project := projectByID(view.Projects, e.AggregateRootID)
		part := findPart(project.Parts, e.PartID, e.ColourID)
		part.Inventory += e.Quantity

	case *ProjectInventoryRemoved:
		project := projectByID(view.Projects, e.AggregateRootID)
		part := findPart(project.Parts, e.PartID, e.ColourID)
		part.Inventory -= e.Quantity

	case *KitCreated:
		kit := parseKitParts(e.Parts)

		view.Kits[e.KitNumber] = kit

		for _, project := range view.Projects {
			calculateKitFulfillment(project, e.KitNumber, kit)
		}

	}

	return view
}

func projectByID(all map[ProjectName]*ProjectView, id uuid.UUID) *ProjectView {
	for _, p := range all {
		if p.ID == id {
			return p
		}
	}
	return nil
}

func findPart(parts []*ProjectPartView, partID LDrawPart, colourID BrickLinkColour) *ProjectPartView {

	for _, part := range parts {
		if part.ID == partID && part.ColourID == colourID {
			return part
		}
	}

	return nil
}

func calculateKitFulfillment(project *ProjectView, kitNumber KitNumber, kitParts map[PartKey]int) {
	fulfilled := map[PartKey]int{}

	for _, part := range project.Parts {

		if quantity, found := kitParts[part.Key]; found {
			fulfilled[part.Key] += quantity
		}
	}

	if len(fulfilled) > 0 {
		project.Kits[kitNumber] = fulfilled
	} else {
		delete(project.Kits, kitNumber)
	}
}

func parseKitParts(parts []Part) map[PartKey]int {

	kp := make(map[PartKey]int, len(parts))

	for _, p := range parts {
		kp[CreatePartKey(p.ID, p.Colour.ID)] = p.Quantity
	}

	return kp
}
