package allkits

import (
	"brickrecon/lego"
)

type PartLoader func(lego.PartKey) *lego.Part

func newPartView(load PartLoader, key lego.PartKey, quantity int) *PartView {

	part := load(key)

	return &PartView{
		Key:        key,
		ID:         part.PartID,
		Name:       part.Name,
		ColourID:   part.ColourID,
		ColourName: part.ColourName,
		ImagePath:  part.ImagePath,
		ColourHex:  part.ColourHex,
		Quantity:   quantity,
	}

}
