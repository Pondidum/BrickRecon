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
		Kits:     map[lego.KitNumber]KitView{},
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
			Kits: map[lego.KitNumber]KitView{},
		}

	case *lego.ProjectPartsAdded:
		project := projectByID(view.Projects, e.AggregateRootID)
		for _, part := range e.Parts {
			project.Parts = append(project.Parts, toProjectPartView(part))
		}

		for _, kit := range view.Kits {
			calculateKitFulfillment(project, kit)
		}

	case *lego.ProjectInventoryAdded:
		project := projectByID(view.Projects, e.AggregateRootID)
		part := findPart(project.Parts, e.PartID, e.ColourID)
		part.Inventory += e.Quantity

	case *lego.ProjectInventoryRemoved:
		project := projectByID(view.Projects, e.AggregateRootID)
		part := findPart(project.Parts, e.PartID, e.ColourID)
		part.Inventory -= e.Quantity

	case *lego.KitAddedToProject:
		project := projectByID(view.Projects, e.AggregateRootID)
		for _, pq := range e.Parts {
			part := findPart(project.Parts, pq.PartID, pq.ColourID)
			part.Inventory += pq.Quantity
		}

	case *lego.KitCreated:
		kit := createKitView(e)

		view.Kits[e.KitNumber] = kit

		for _, project := range view.Projects {
			calculateKitFulfillment(project, kit)
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

func calculateKitFulfillment(project *ProjectView, kit KitView) {
	fulfilled := map[PartKey]int{}
	total := 0

	for _, part := range project.Parts {

		if quantity, found := kit.Parts[part.Key]; found {
			fulfilled[part.Key] += quantity
			total += quantity
		}
	}

	if len(fulfilled) > 0 {
		project.Kits[kit.Number] = KitView{
			Number:     kit.Number,
			Name:       kit.Name,
			Parts:      fulfilled,
			TotalParts: total,
		}
	} else {
		delete(project.Kits, kit.Number)
	}
}

func createKitView(event *lego.KitCreated) KitView {

	kp := make(map[PartKey]int, len(event.Parts))

	for _, p := range event.Parts {
		kp[CreatePartKey(p.ID, p.Colour.ID)] = p.Quantity
	}

	return KitView{
		Number: event.KitNumber,
		Name:   event.KitName,
		Parts:  kp,
	}
}
