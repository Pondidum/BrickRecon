package all_kits

import (
	"brickrecon/eventstore"
	"brickrecon/lego"
)

type AllKitsView struct {
	Kits map[lego.KitNumber]*KitView
}

type KitView struct {
	ID     eventstore.AggregateID
	Name   lego.KitName
	Number lego.KitNumber

	Parts []PartView
}

type PartView struct {
	Key        lego.PartKey
	ID         lego.LDrawPart
	Name       lego.PartName
	ColourID   lego.BrickLinkColour
	ColourName lego.ColourName
	ColourHex  lego.HexColour

	Quantity int
}
