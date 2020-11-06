package app

import (
	"brickrecon/lego"
	"brickrecon/preen"
	"strconv"
)

type filter struct {
	HasColour bool
	Colour    lego.BrickLinkColour

	HasInventory bool
	Inventory    string
}

func createFilter(pc *preen.PreenContext) filter {

	colourRaw := pc.QueryValue("filter-colour")
	colour, err := strconv.Atoi(colourRaw)

	inventoryRaw := pc.QueryValue("filter-inventory")

	return filter{
		HasColour: err == nil,
		Colour:    lego.BrickLinkColour(colour),

		HasInventory: inventoryRaw != "" && inventoryRaw != "all",
		Inventory:    inventoryRaw,
	}
}

func (f *filter) Parts(parts []*PartWithKitPart) []*PartWithKitPart {

	result := []*PartWithKitPart{}

	for _, part := range parts {

		colourMatch := f.isColourMatch(part)
		inventoryMatch := f.isInventoryMatch(part)

		if colourMatch && inventoryMatch {
			result = append(result, part)
		}
	}

	return result
}

func (f *filter) isColourMatch(part *PartWithKitPart) bool {
	return f.HasColour == false || part.ColourID == f.Colour
}

func (f *filter) isInventoryMatch(part *PartWithKitPart) bool {

	if f.HasInventory == false {
		return true
	}

	if f.Inventory == "needed" {
		return part.Inventory < part.Quantity
	}

	if f.Inventory == "owned" {
		return part.Inventory >= part.Quantity
	}

	return true
}
