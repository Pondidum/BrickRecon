package all_projects

import (
	"brickrecon/eventstore"
	"brickrecon/lego"
	"fmt"
	"reflect"

	"github.com/mitchellh/mapstructure"
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
		Key:        lego.CreatePartKey(part.ID, part.Colour.ID),
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

	project := projectByID(view.Projects, event.Meta().AggregateRootID)

	switch e := event.(type) {

	case *lego.ProjectCreated:
		project = &ProjectView{
			ID:      e.AggregateRootID,
			Name:    e.Name,
			Kits:    map[lego.KitNumber]KitView{},
			Colours: []*ColourView{},
		}

		audit(project, event, "Project created")

		view.Names = append(view.Names, e.Name)
		view.Projects[e.Name] = project

	case *lego.ProjectPartsAdded:
		for _, part := range e.Parts {
			project.Parts = append(project.Parts, toProjectPartView(part))
			project.Colours = appendNewColours(project.Colours, part)
		}

		for _, kit := range view.Kits {
			calculateKitFulfillment(project, kit)
		}

		audit(project, e, "%v parts added", len(e.Parts))

	case *lego.ProjectInventoryAdded:
		part := findPart(project.Parts, e.PartID, e.ColourID)
		part.Inventory += e.Quantity

		audit(project, e, "Added %v %s %s (%s)", e.Quantity, part.ColourName, part.Name, part.ID)

	case *lego.ProjectInventoryRemoved:
		part := findPart(project.Parts, e.PartID, e.ColourID)
		part.Inventory -= e.Quantity

		audit(project, e, "Removed %v %s %s (%s)", e.Quantity, part.ColourName, part.Name, part.ID)

	case *lego.KitAddedToProject:
		for _, pq := range e.Parts {
			part := findPart(project.Parts, pq.PartID, pq.ColourID)
			part.Inventory += pq.Quantity
		}
		audit(project, e, "%s (Kit %s) applied", e.KitName, e.KitNumber)

	case *lego.KitCreated:
		kit := createKitView(e)

		view.Kits[e.KitNumber] = kit

		for _, project := range view.Projects {
			calculateKitFulfillment(project, kit)
		}

	case *lego.WantedListExported:
		project.BrickLinkXml = e.Markup

		audit(project, e, "WantedList XML generated")
	}

	return view
}

func appendNewColours(unique []*ColourView, part lego.Part) []*ColourView {

	for _, view := range unique {
		if view.ID == part.Colour.ID {
			return unique
		}
	}

	unique = append(unique, &ColourView{
		ID:   part.Colour.ID,
		Name: part.Colour.Name,
		Hex:  part.Colour.Hex,
	})

	return unique
}

func audit(project *ProjectView, event eventstore.Event, format string, args ...interface{}) {

	desc := &EventDescription{
		Timestamp:   event.Meta().Timestamp,
		Type:        event.Meta().Type,
		Description: fmt.Sprintf(format, args...),
		Additional:  map[string]interface{}{},
	}

	decoder, _ := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		DecodeHook: func(src reflect.Type, dest reflect.Type, in interface{}) (interface{}, error) {
			return in, nil
		},
		Result: &desc.Additional,
	})

	decoder.Decode(event)
	// mapstructure.Decode(event, &desc.Additional)
	delete(desc.Additional, "EventMeta")
	delete(desc.Additional, "ID")

	project.Events = append(project.Events, desc)
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

	for _, p := range event.Parts {
		kp[lego.CreatePartKey(p.ID, p.Colour.ID)] = p.Quantity
	}

	return KitView{
		Number: event.KitNumber,
		Name:   event.KitName,
		Parts:  kp,
	}
}
