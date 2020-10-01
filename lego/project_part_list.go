package lego

import "fmt"

type ProjectPartList struct {
	parts map[PartKey]*ProjectPart
}

type ProjectPart struct {
	Part

	Inventory int
}

func (p *ProjectPart) HasSpares() bool {
	return p.Inventory > p.Quantity
}

func (p *ProjectPart) IsFulfilled() bool {
	return p.Inventory >= p.Quantity
}

func (p *ProjectPart) NeededQuantity() int {
	return p.Quantity - p.Inventory
}

func NewPartsList() *ProjectPartList {
	list := ProjectPartList{
		parts: map[PartKey]*ProjectPart{},
	}

	return &list
}

func (m *ProjectPartList) Add(part Part) {

	key := CreatePartKey(part.ID, part.Colour.ID)
	existing, found := m.FindPartByKey(key)

	if found {
		existing.Quantity += part.Quantity
		return
	}

	m.parts[key] = &ProjectPart{Part: part, Inventory: 0}
}

func (m *ProjectPartList) AddInventory(partID LDrawPart, colourID BrickLinkColour, quantity int) error {

	part, found := m.FindPart(partID, colourID)

	if !found {
		return fmt.Errorf("No part with id %s and colour %v found", partID, colourID)
	}

	part.Inventory += quantity

	if part.Inventory < 0 {
		part.Inventory = 0
	}

	return nil
}

func (m *ProjectPartList) FindPart(partID LDrawPart, colourID BrickLinkColour) (*ProjectPart, bool) {

	for _, p := range m.parts {

		if p.ID == partID && p.Colour.ID == colourID {
			return p, true
		}
	}

	return nil, false
}

func (m *ProjectPartList) FindPartByKey(partKey PartKey) (*ProjectPart, bool) {
	part, found := m.parts[partKey]

	return part, found
}
