package lego

import (
	"brickrecon/eventstore"

	uuid "github.com/satori/go.uuid"
)

type AllKitsView struct {
	Kits map[string]*KitView
}

type KitView struct {
	ID     uuid.UUID
	Name   string
	Number string

	Parts []PartView
}

type PartView struct {
	ID         PartID
	Name       string
	ColourID   BrickLinkColour
	ColourName string

	Quantity int
}

func toPartView(parts []Part) []PartView {

	views := make([]PartView, len(parts))

	for i, part := range parts {
		views[i] = PartView{
			ID:         part.ID,
			Name:       part.Name,
			ColourID:   part.Colour.ID,
			ColourName: part.Colour.Name,
			Quantity:   part.Quantity,
		}
	}

	return views
}

func KitsInitialState() interface{} {
	return &AllKitsView{
		Kits: map[string]*KitView{},
	}
}

func KitsProjector(state interface{}, event eventstore.Event) interface{} {
	view := state.(*AllKitsView)

	switch e := event.(type) {

	case *KitCreated:
		view.Kits[e.KitNumber] = &KitView{
			ID:     e.AggregateID(),
			Name:   e.KitName,
			Number: e.KitNumber,
			Parts:  toPartView(e.Parts),
		}
	}

	return view
}
