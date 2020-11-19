package allkits

import (
	"brickrecon/lego"
)

type PartLoader func(lego.PartKey) *lego.PartA

func newPartView(load PartLoader, key lego.PartKey, quantity int) *PartView {

	part := load(key)
	hex := lego.HexColour("")

	if colour, found := lego.LookupColourLDraw(part.ColourID); found {
		hex = colour.Hex
	}

	return &PartView{
		Key:        key,
		ID:         part.PartID,
		Name:       part.Name,
		ColourID:   part.ColourID,
		ColourName: part.ColourName,
		ImagePath:  part.ImagePath,
		Quantity:   quantity,
		ColourHex:  hex,
	}

}
