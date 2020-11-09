package all_projects

import (
	"brickrecon/eventstore"
	"brickrecon/lego"
	"fmt"
	"reflect"
	"sort"

	"github.com/mitchellh/mapstructure"
)

func toProjectPartView(part *lego.Part) *ProjectPartView {
	return &ProjectPartView{
		ID:         part.Aliases.LDrawID,
		Name:       part.Name,
		ColourID:   part.Colour.ID,
		ColourName: part.Colour.Name,
		ColourHex:  part.Colour.Hex,
		Quantity:   part.Quantity,
		Key:        part.Key,
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
			Stats:   &ProjectStatsView{},
		}

		audit(project, event, "Project created")

		view.Names = append(view.Names, e.Name)
		view.Projects[e.Name] = project

	case *lego.ProjectPartsAdded:
		addParts(project, e.Parts)

		for _, kit := range view.Kits {
			calculateKitFulfillment(project, kit)
		}

		calculateStats(project)
		audit(project, e, "%v parts added", len(e.Parts))

	case *lego.PartsChanged:
		addParts(project, e.Additions)
		removeParts(project, e.Removals)

		for _, kit := range view.Kits {
			calculateKitFulfillment(project, kit)
		}

		calculateStats(project)
		audit(project, e, "%v parts changed", len(e.Additions)+len(e.Removals))

	case *lego.ProjectInventoryAdded:
		part, _ := findPart(project.Parts, e.PartID, e.ColourID)
		part.Inventory += e.Quantity

		calculateStats(project)
		audit(project, e, "Added %v %s %s (%s)", e.Quantity, part.ColourName, part.Name, part.ID)

	case *lego.ProjectInventoryRemoved:
		part, _ := findPart(project.Parts, e.PartID, e.ColourID)
		part.Inventory -= e.Quantity

		calculateStats(project)
		audit(project, e, "Removed %v %s %s (%s)", e.Quantity, part.ColourName, part.Name, part.ID)

	case *lego.KitAddedToProject:
		for _, pq := range e.Parts {
			part, _ := findPart(project.Parts, pq.PartID, pq.ColourID)
			part.Inventory += pq.Quantity
		}

		calculateStats(project)
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

func addParts(project *ProjectView, parts []*lego.Part) {
	for _, part := range parts {
		project.Parts = append(project.Parts, toProjectPartView(part))
		project.Colours = appendNewColours(project.Colours, part)
	}
}

func removeParts(project *ProjectView, parts map[lego.PartKey]int) {

	for key, amount := range parts {
		partID, colourID := lego.ParsePartKey(key)
		part, index := findPart(project.Parts, partID, colourID)

		part.Quantity -= amount

		if part.Quantity <= 0 {
			project.Parts = append(project.Parts[:index], project.Parts[index+1:]...)
		}
	}
}

func appendNewColours(unique []*ColourView, part *lego.Part) []*ColourView {

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

	sort.Slice(unique, func(i, j int) bool {
		return unique[i].Name < unique[j].Name
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

func calculateStats(project *ProjectView) {
	totalQuantity := 0
	totalInventory := 0

	for _, part := range project.Parts {
		totalQuantity += part.Quantity
		totalInventory += part.Inventory
	}

	project.Stats.TotalPartsQuantity = totalQuantity
	project.Stats.TotalPartsInventory = totalInventory
	project.Stats.PercentComplete = int(float64(totalInventory) / float64(totalQuantity) * 100)
}

func projectByID(all map[lego.ProjectName]*ProjectView, id eventstore.AggregateID) *ProjectView {
	for _, p := range all {
		if p.ID == id {
			return p
		}
	}
	return nil
}

func findPart(parts []*ProjectPartView, partID lego.LDrawPart, colourID lego.BrickLinkColour) (*ProjectPartView, int) {

	for i, part := range parts {
		if part.ID == partID && part.ColourID == colourID {
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

	for _, p := range event.Parts {
		kp[p.Key] = p.Quantity
	}

	return KitView{
		Number: event.KitNumber,
		Name:   event.KitName,
		Parts:  kp,
	}
}
