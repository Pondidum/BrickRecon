package lego

import (
	"brickrecon/eventstore"

	uuid "github.com/satori/go.uuid"
)

type AllKitsView struct {
	Kits map[KitNumber]*KitView
}

type KitView struct {
	ID     uuid.UUID
	Name   KitName
	Number KitNumber

	Parts []PartView
}

type PartView struct {
	ID         LDrawPart
	Name       PartName
	ColourID   BrickLinkColour
	ColourName ColourName
	ColourHex  HexColour

	Quantity int
}

var KitsProjectionName string = "kits"

type KitsProjection struct{}

func (p *KitsProjection) Name() string {
	return KitsProjectionName
}

func (p *KitsProjection) CreateState() interface{} {
	return &AllKitsView{
		Kits: map[KitNumber]*KitView{},
	}
}

func (p *KitsProjection) Project(state interface{}, event eventstore.Event) interface{} {
	view := state.(*AllKitsView)

	switch e := event.(type) {

	case *KitCreated:
		view.Kits[e.KitNumber] = &KitView{
			ID:     e.AggregateRootID,
			Name:   e.KitName,
			Number: e.KitNumber,
			Parts:  toPartView(e.Parts),
		}
	}

	return view
}

func toPartView(parts []Part) []PartView {

	views := make([]PartView, len(parts))

	for i, part := range parts {
		views[i] = PartView{
			ID:         part.ID,
			Name:       part.Name,
			ColourID:   part.Colour.ID,
			ColourName: part.Colour.Name,
			ColourHex:  part.Colour.Hex,
			Quantity:   part.Quantity,
		}
	}

	return views
}
