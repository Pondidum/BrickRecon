package all_projects

import (
	"brickrecon/eventstore"
	"brickrecon/lego"

	uuid "github.com/satori/go.uuid"
)

func toProjectPartView(part lego.Part) *ProjectPartView {
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

var ProjectionName string = "projects"

type ProjectsProjection struct{}

func (p *ProjectsProjection) Name() string {
	return ProjectionName
}

func (p *ProjectsProjection) CreateState() interface{} {
	return &AllProjectsView{
		Names:    []lego.ProjectName{},
		Projects: map[lego.ProjectName]*ProjectView{},
		Kits:     map[lego.KitNumber]map[PartKey]int{},
	}
}

func (p *ProjectsProjection) Project(state interface{}, event eventstore.Event) interface{} {
	view := state.(*AllProjectsView)

	switch e := event.(type) {

	case *lego.ProjectCreated:
		view.Names = append(view.Names, e.Name)
		view.Projects[e.Name] = &ProjectView{
			ID:   e.AggregateRootID,
			Name: e.Name,
			Kits: map[lego.KitNumber]map[PartKey]int{},
		}

	case *lego.ProjectPartsAdded:
		project := projectByID(view.Projects, e.AggregateRootID)
		for _, part := range e.Parts {
			project.Parts = append(project.Parts, toProjectPartView(part))
		}

		for kn, kit := range view.Kits {
			calculateKitFulfillment(project, kn, kit)
		}

	case *lego.ProjectInventoryAdded:
		project := projectByID(view.Projects, e.AggregateRootID)
		part := findPart(project.Parts, e.PartID, e.ColourID)
		part.Inventory += e.Quantity

	case *lego.ProjectInventoryRemoved:
		project := projectByID(view.Projects, e.AggregateRootID)
		part := findPart(project.Parts, e.PartID, e.ColourID)
		part.Inventory -= e.Quantity

	case *lego.KitCreated:
		kit := parseKitParts(e.Parts)

		view.Kits[e.KitNumber] = kit

		for _, project := range view.Projects {
			calculateKitFulfillment(project, e.KitNumber, kit)
		}

	}

	return view
}

func projectByID(all map[lego.ProjectName]*ProjectView, id uuid.UUID) *ProjectView {
	for _, p := range all {
		if p.ID == id {
			return p
		}
	}
	return nil
}

func findPart(parts []*ProjectPartView, partID lego.LDrawPart, colourID lego.BrickLinkColour) *ProjectPartView {

	for _, part := range parts {
		if part.ID == partID && part.ColourID == colourID {
			return part
		}
	}

	return nil
}

func calculateKitFulfillment(project *ProjectView, kitNumber lego.KitNumber, kitParts map[PartKey]int) {
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

func parseKitParts(parts []lego.Part) map[PartKey]int {

	kp := make(map[PartKey]int, len(parts))

	for _, p := range parts {
		kp[CreatePartKey(p.ID, p.Colour.ID)] = p.Quantity
	}

	return kp
}
