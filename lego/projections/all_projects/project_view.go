package all_projects

import (
	"brickrecon/eventstore"
	"brickrecon/lego"
	"fmt"
	"reflect"
	"sort"
	"time"

	"github.com/mitchellh/mapstructure"
)

type ProjectView struct {
	ID   eventstore.AggregateID
	Name lego.ProjectName

	Parts   []*ProjectPartView
	Kits    map[lego.KitNumber]KitView
	Colours []*ColourView

	Stats *ProjectStatsView

	BrickLinkXml string

	Events []*EventDescription
}

type EventDescription struct {
	Timestamp   time.Time
	Type        string
	Description string
	Additional  map[string]interface{}
}

func newProjectView(id eventstore.AggregateID, name lego.ProjectName) *ProjectView {
	return &ProjectView{
		ID:      id,
		Name:    name,
		Kits:    map[lego.KitNumber]KitView{},
		Colours: []*ColourView{},
		Stats:   &ProjectStatsView{},
	}
}

func (project *ProjectView) addParts(load PartLoader, parts []*lego.Part) {
	for _, part := range parts {
		view := newPartView(load, part.Key, part.Quantity)
		project.Parts = append(project.Parts, view)
		project.Colours = appendNewColours(project.Colours, part)
	}
}

func (project *ProjectView) removeParts(parts map[lego.PartKey]int) {

	for key, amount := range parts {
		part, index := findPart(project.Parts, key)

		part.Quantity -= amount

		if part.Quantity <= 0 {
			project.Parts = append(project.Parts[:index], project.Parts[index+1:]...)
		}
	}
}

func (project *ProjectView) calculateStats() {
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

func appendNewColours(unique []*ColourView, part *lego.Part) []*ColourView {

	for _, view := range unique {
		if view.ID == part.Colour.Aliases.LDrawID {
			return unique
		}
	}

	unique = append(unique, &ColourView{
		ID:   part.Colour.Aliases.LDrawID,
		Name: part.Colour.Name,
		Hex:  part.Colour.Hex,
	})

	sort.Slice(unique, func(i, j int) bool {
		return unique[i].Name < unique[j].Name
	})
	return unique
}

func (project *ProjectView) audit(event eventstore.Event, format string, args ...interface{}) {

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
