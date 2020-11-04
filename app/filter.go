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

func filterParts(parts []*PartWithKitPart, f filter) []*PartWithKitPart {

	result := []*PartWithKitPart{}

	for _, part := range parts {

		colourMatch := f.HasColour == false || part.ColourID == f.Colour

		inventoryMatch := f.HasInventory == false || isInventoryMatch(part, f.Inventory)

		if colourMatch && inventoryMatch {
			result = append(result, part)
		}
	}

	return result
}

func isInventoryMatch(part *PartWithKitPart, inventory string) bool {

	if inventory == "needed" {
		return part.Inventory < part.Quantity
	}

	if inventory == "owned" {
		return part.Inventory >= part.Quantity
	}

	return true
}
