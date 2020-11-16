package all_projects

import (
	"brickrecon/lego"
)

type PartLoader func(lego.PartKey) *lego.PartAggregate

func newPartView(load PartLoader, key lego.PartKey, quantity int) *ProjectPartView {

	part := load(key)
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
