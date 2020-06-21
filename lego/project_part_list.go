package lego

import "fmt"

type ProjectPartList struct {
	parts []*ProjectPart
}

type ProjectPart struct {
	Part

	Inventory int
}

func (p *ProjectPart) HasSpares() bool {
	return p.Inventory > p.Quantity
}

func NewPartsList(parts []Part) *ProjectPartList {
	list := ProjectPartList{
		parts: make([]*ProjectPart, len(parts)),
	}

	for i, p := range parts {
		list.parts[i] = &ProjectPart{Part: p, Inventory: 0}
	}

	return &list
}

func (m *ProjectPartList) Add(part Part) {

	id := part.ID
	colour := part.Colour.ID

	existing, found := m.FindPart(id, colour)

	if found {
		existing.Quantity += part.Quantity
		return
	}

	m.parts = append(m.parts, &ProjectPart{Part: part, Inventory: 0})
}

func (m *ProjectPartList) AddInventory(partID PartID, colourID int, quantity int) error {

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

func (m *ProjectPartList) FindPart(partID PartID, colourID int) (*ProjectPart, bool) {

	for _, p := range m.parts {

		if p.ID == partID && p.Colour.ID == colourID {
			return p, true
		}
	}

	return nil, false
}
