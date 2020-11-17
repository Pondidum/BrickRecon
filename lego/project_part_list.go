package lego

import "fmt"

type ProjectPartList struct {
	parts map[PartKey]*ProjectPart
}

type ProjectPart struct {
	Quantity  int
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

func (l *ProjectPartList) All() map[PartKey]*ProjectPart {
	return l.parts
}

func (m *ProjectPartList) Add(part *Part) {
	m.AddPart(part.Key, part.Quantity)
}

func (m *ProjectPartList) AddPart(key PartKey, quantity int) {
	existing, found := m.FindPart(key)

	if found {
		existing.Quantity += quantity
		return
	}

	m.parts[key] = &ProjectPart{Quantity: quantity, Inventory: 0}
}

func (m *ProjectPartList) Remove(key PartKey, quantity int) {

	if part, found := m.FindPart(key); found {
		part.Quantity -= quantity

		if part.Quantity <= 0 {
			delete(m.parts, key)
		}
	}

}

func (m *ProjectPartList) AddInventory(key PartKey, quantity int) error {

	part, found := m.FindPart(key)

	if !found {
		return fmt.Errorf("No part with id %s found", key)
	}

	part.Inventory += quantity

	if part.Inventory < 0 {
		part.Inventory = 0
	}

	return nil
}

func (m *ProjectPartList) FindPart(partKey PartKey) (*ProjectPart, bool) {
	part, found := m.parts[partKey]

	return part, found
}

func (m *ProjectPartList) Diff(other *ProjectPartList) map[PartKey]int {

	deltas := map[PartKey]int{}

	for key, op := range other.parts {
		if p, found := m.parts[key]; found {

			quantityChange := op.Quantity - p.Quantity

			if quantityChange != 0 {
				deltas[key] = quantityChange
			}
		} else {
			deltas[key] = op.Quantity
		}
	}

	for key, p := range m.parts {

		if _, found := other.parts[key]; !found {
			deltas[key] = p.Quantity * -1
		}
	}

	return deltas
}
