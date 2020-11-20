package allkits

import (
	"brickrecon/lego"
	"context"
)

type PartLoader func(ctx context.Context, key lego.PartKey) *lego.Part

func newPartView(ctx context.Context, load PartLoader, key lego.PartKey, quantity int) *PartView {

	part := load(ctx, key)

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
