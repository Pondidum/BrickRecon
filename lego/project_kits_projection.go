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

type ProjectKitsView struct {
	Kits     map[KitNumber]map[PartKey]int
	Projects map[uuid.UUID]*ProjectKit
}

type ProjectKit struct {
	Parts []PartRequirement
	Kits  map[KitNumber]map[PartKey]int
}

type PartRequirement struct {
	PartID   LDrawPart
	Colour   BrickLinkColour
	Key      PartKey
	Quantity int
}

var ProjectKitsProjectionName string = "project_kits"

type ProjectKitsProjection struct{}

func (p *ProjectKitsProjection) Name() string {
	return ProjectKitsProjectionName
}

func (p *ProjectKitsProjection) CreateState() interface{} {
	return &ProjectKitsView{
		Kits:     map[KitNumber]map[PartKey]int{},
		Projects: map[uuid.UUID]*ProjectKit{},
	}
}

func (p *ProjectKitsProjection) Project(state interface{}, event eventstore.Event) interface{} {
	view := state.(*ProjectKitsView)

	switch e := event.(type) {

	case *ProjectPartsAdded:
		project := &ProjectKit{
			Parts: parseRequirements(e.Parts),
			Kits:  map[KitNumber]map[PartKey]int{},
		}

		view.Projects[e.AggregateRootID] = project

		for kn, kit := range view.Kits {
			fill(project, kn, kit)
		}

	case *KitCreated:
		kit := parseKitParts(e.Parts)

		view.Kits[e.KitNumber] = kit

		for _, project := range view.Projects {
			fill(project, e.KitNumber, kit)
		}

	}

	return view
}

func fill(project *ProjectKit, kitNumber KitNumber, kitParts map[PartKey]int) {
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

func parseRequirements(parts []Part) []PartRequirement {
	req := make([]PartRequirement, len(parts))

	for i, p := range parts {
		req[i] = PartRequirement{
			PartID:   p.ID,
			Colour:   p.Colour.ID,
			Key:      CreatePartKey(p.ID, p.Colour.ID),
			Quantity: p.Quantity,
		}
	}

	return req
}

func parseKitParts(parts []Part) map[PartKey]int {

	kp := make(map[PartKey]int, len(parts))

	for _, p := range parts {
		kp[CreatePartKey(p.ID, p.Colour.ID)] = p.Quantity
	}

	return kp
}
