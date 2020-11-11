package all_projects

import (
	"brickrecon/lego"
	"fmt"
)

type AllProjectsView struct {
	Names    []lego.ProjectName
	Projects map[lego.ProjectName]*ProjectView
	Kits     map[lego.KitNumber]KitView
}

type ProjectPartView struct {
	Key        lego.PartKey
	ID         lego.LDrawPart
	Name       lego.PartName
	ColourID   lego.LDrawColour
	ColourName lego.ColourName
	ColourHex  lego.HexColour

	ImagePath string

	Quantity  int
	Inventory int
}

func newPartViewFromPart(part *lego.Part) *ProjectPartView {
	return &ProjectPartView{
		ID:         part.Aliases.LDrawID,
		Name:       part.Name,
		ColourID:   part.Colour.Aliases.LDrawID,
		ColourName: part.Colour.Name,
		ColourHex:  part.Colour.Hex,
		ImagePath:  fmt.Sprintf("%s-%v.png", part.Aliases.BrickLinkID, part.Colour.Aliases.BrickLinkID),
		Quantity:   part.Quantity,
		Key:        part.Key,
	}
}

type KitView struct {
	Number lego.KitNumber
	Name   lego.KitName

	Parts      map[lego.PartKey]int
	TotalParts int
}

type ColourView struct {
	ID   lego.LDrawColour
	Name lego.ColourName
	Hex  lego.HexColour
}

type ProjectStatsView struct {
	TotalPartsQuantity  int
	TotalPartsInventory int
	PercentComplete     int
}
