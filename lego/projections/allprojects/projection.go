package allprojects

import (
	"brickrecon/eventstore"
	"brickrecon/lego"
	"context"
)

var ProjectionName string = "projects"

func NewProjectsProjection(es eventstore.EventStore) *projectsProjection {
	return &projectsProjection{
		partLoader: func(key lego.PartKey) *lego.Part {
			part := lego.BlankPart()
			es.LoadAggregate(context.Background(), eventstore.AggregateID(key), part)
			return part
		},
	}
}

type projectsProjection struct {
	partLoader PartLoader
}

func (p *projectsProjection) Name() string {
	return ProjectionName
}

func (p *projectsProjection) CreateState() interface{} {
	return &AllProjectsView{
		Names:    []lego.ProjectName{},
		Projects: map[lego.ProjectName]*ProjectView{},
		Kits:     map[lego.KitNumber]KitView{},
	}
}

func (p *projectsProjection) Project(state interface{}, event eventstore.Event) interface{} {
	view := state.(*AllProjectsView)

	project := projectByID(view.Projects, event.Meta().AggregateRootID)

	switch e := event.(type) {

	case *lego.ProjectCreated:
		project = newProjectView(e.AggregateRootID, e.Name)
		project.audit(event, "Project created")

		view.Names = append(view.Names, e.Name)
		view.Projects[e.Name] = project

	case *lego.ProjectPartsAdded:
		project.addParts(p.partLoader, e.Parts)

		for _, kit := range view.Kits {
			calculateKitFulfillment(project, kit)
		}

		project.audit(e, "%v parts added", len(e.Parts))

	case *lego.PartsChanged:
		project.addParts(p.partLoader, e.Additions)
		project.removeParts(e.Removals)

		for _, kit := range view.Kits {
			calculateKitFulfillment(project, kit)
		}

		project.audit(e, "%v parts changed", len(e.Additions)+len(e.Removals))

	case *lego.ProjectInventoryAdded:
		part, _ := findPart(project.Parts, e.Part)
		part.Inventory += e.Quantity

		project.audit(e, "Added %v %s %s (%s)", e.Quantity, part.ColourName, part.Name, part.ID)

	case *lego.ProjectInventoryRemoved:
		part, _ := findPart(project.Parts, e.Part)
		part.Inventory -= e.Quantity

		project.audit(e, "Removed %v %s %s (%s)", e.Quantity, part.ColourName, part.Name, part.ID)

	case *lego.KitAddedToProject:
		for key, quantity := range e.Parts {
			part, _ := findPart(project.Parts, key)
			part.Inventory += quantity
		}

		project.audit(e, "%s (Kit %s) applied", e.KitName, e.KitNumber)

	case *lego.KitCreated:
		kit := createKitView(e)

		view.Kits[e.KitNumber] = kit

		for _, project := range view.Projects {
			calculateKitFulfillment(project, kit)
		}

	}

	if project != nil {
		project.calculateStats()
	}

	return view
}

func projectByID(all map[lego.ProjectName]*ProjectView, id eventstore.AggregateID) *ProjectView {
	for _, p := range all {
		if p.ID == id {
			return p
		}
	}
	return nil
}

func findPart(parts []*ProjectPartView, key lego.PartKey) (*ProjectPartView, int) {

	for i, part := range parts {
		if part.Key == key {
			return part, i
		}
	}

	return nil, -1
}

func calculateKitFulfillment(project *ProjectView, kit KitView) {
	fulfilled := map[lego.PartKey]int{}
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

	kp := make(map[lego.PartKey]int, len(event.Parts))

	for key, quantity := range event.Parts {
		kp[key] = quantity
	}

	return KitView{
		Number: event.KitNumber,
		Name:   event.KitName,
		Parts:  kp,
	}
}
