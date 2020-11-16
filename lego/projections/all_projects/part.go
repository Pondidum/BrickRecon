package all_projects

import (
	"brickrecon/eventstore"
	"brickrecon/lego"
	"context"
)

func newPartView(es eventstore.EventStore, key lego.PartKey, quantity int) *ProjectPartView {

	part := lego.BlankPart()
	es.LoadAggregate(context.Background(), eventstore.AggregateID(key), part)

	hex := lego.HexColour("")

	if colour, found := lego.LookupColourLDraw(part.Colour); found {
		hex = colour.Hex
	}

	return &ProjectPartView{
		Key:        key,
		ID:         part.Number,
		Name:       part.Name,
		ColourID:   part.Colour,
		ColourName: part.ColourName,
		ImagePath:  part.ImagePath,
		Quantity:   quantity,
		ColourHex:  hex,
	}

}
