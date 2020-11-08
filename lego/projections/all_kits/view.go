package all_kits

import (
	"brickrecon/lego"

	uuid "github.com/satori/go.uuid"
)

type AllKitsView struct {
	Kits map[lego.KitNumber]*KitView
}

type KitView struct {
	ID     uuid.UUID
	Name   lego.KitName
	Number lego.KitNumber

	Parts []PartView
}

type PartView struct {
	Key        lego.PartKey
	Name       lego.PartName
	ColourID   lego.BrickLinkColour
	ColourName lego.ColourName
	ColourHex  lego.HexColour

	Quantity int
}
