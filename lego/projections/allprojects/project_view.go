package allprojects

import (
	"brickrecon/eventstore"
	"brickrecon/lego"
	"context"
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

func (project *ProjectView) addParts(ctx context.Context, load PartLoader, parts map[lego.PartKey]int) {
	for key, quantity := range parts {
		view := newPartView(ctx, load, key, quantity)
		project.Parts = append(project.Parts, view)
		project.Colours = appendNewColours(project.Colours, view)
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

func appendNewColours(unique []*ColourView, part *ProjectPartView) []*ColourView {

	for _, view := range unique {
		if view.ID == part.ColourID {
			return unique
		}
	}

	unique = append(unique, &ColourView{
		ID:   part.ColourID,
		Name: part.ColourName,
		Hex:  part.ColourHex,
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
