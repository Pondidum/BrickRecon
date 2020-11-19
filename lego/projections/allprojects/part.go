package allprojects

import (
	"brickrecon/lego"
)

type PartLoader func(lego.PartKey) *lego.Part

func newPartView(load PartLoader, key lego.PartKey, quantity int) *ProjectPartView {

	part := load(key)

	return &ProjectPartView{
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
