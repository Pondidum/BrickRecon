package lego

import (
	"brickrecon/eventstore"
	"errors"
	"fmt"
)

type ProjectName string

type Project struct {
	*eventstore.Aggregator

	Name  ProjectName
	parts *ProjectPartList
}

func BlankProject() *Project {
	project := Project{
		parts: NewPartsList(),
	}
	project.Aggregator = eventstore.NewAggregator(project.on)

	return &project
}

func NewProject(name ProjectName, parts map[PartKey]int) *Project {
	project := BlankProject()
	project.Apply(&ProjectCreated{ID: eventstore.NewAggregateID(), Name: name})
	project.Apply(&ProjectPartsAdded{
		Parts: parts,
	})

	return project
}

func (prj *Project) FindPart(part PartKey) (*ProjectPart, bool) {
	return prj.parts.FindPart(part)
}

func (prj *Project) AddInventory(part PartKey, quantity int) error {

	if quantity <= 0 {
		return errors.New("Quantity must be greater than 0")
	}

	if _, found := prj.parts.FindPart(part); !found {
		return fmt.Errorf("No part with id %s found", part)
	}

	prj.Apply(&ProjectInventoryAdded{
		EventMeta: eventstore.EventMeta{EventVersion: 1},
		Part:      part,
		Quantity:  quantity,
	})

	return nil
}

func (prj *Project) RemoveInventory(part PartKey, quantity int) error {

	if quantity <= 0 {
		return errors.New("Quantity must be greater than 0")
	}

	if _, found := prj.parts.FindPart(part); !found {
		return fmt.Errorf("No part with id %s found", part)
	}

	prj.Apply(&ProjectInventoryRemoved{
		EventMeta: eventstore.EventMeta{EventVersion: 1},
		Part:      part,
		Quantity:  quantity,
	})

	return nil
}

func (prj *Project) UpdateInventory(inventoryState map[PartKey]int) error {

	for key, quantity := range inventoryState {

		partID, colourID := ParsePartKey(key)

		current, found := prj.parts.FindPart(key)
		if !found {
			return fmt.Errorf("No part with id %s and colour %v found", partID, colourID)
		}

		diff := quantity - current.Inventory

		if diff < 0 {
			prj.Apply(&ProjectInventoryRemoved{
				EventMeta: eventstore.EventMeta{EventVersion: 1},
				Part:      key,
				Quantity:  diff * -1,
			})
		}

		if diff > 0 {
			prj.Apply(&ProjectInventoryAdded{
				EventMeta: eventstore.EventMeta{EventVersion: 1},
				Part:      key,
				Quantity:  diff,
			})
		}
	}

	return nil
}

func (prj *Project) AddKitContents(number KitNumber, name KitName, parts map[PartKey]int) {

	if len(parts) == 0 {
		return
	}

	prj.Apply((&KitAddedToProject{KitNumber: number, KitName: name, Parts: parts}))
}

func (prj *Project) ReplaceParts(parts map[PartKey]int) map[PartKey]int {

	changes := prj.parts.Diff(parts)

	if len(changes) == 0 {
		return changes
	}

	event := &PartsChanged{
		EventMeta: eventstore.EventMeta{EventVersion: 1},
		Additions: map[PartKey]int{},
		Removals:  map[PartKey]int{},
	}

	for key, change := range changes {
		if change < 0 {
			event.Removals[key] = change * -1
		}

		if change > 0 {
			event.Additions[key] = change
		}
	}

	prj.Apply(event)

	return changes
}

func (prj *Project) Parts() map[PartKey]*ProjectPart {
	return prj.parts.All()
}

func (prj *Project) Diff(parts map[PartKey]int) map[PartKey]int {
	return prj.parts.Diff(parts)
}

func (prj *Project) on(event eventstore.Event) {

	switch e := event.(type) {

	case *ProjectCreated:
		prj.SetID(e.ID)
		prj.Name = e.Name

	case *ProjectPartsAdded:
		for key, quantity := range e.Parts {
			prj.parts.Add(key, quantity)
		}

	case *ProjectInventoryAdded:
		prj.parts.AddInventory(e.Part, e.Quantity)

	case *ProjectInventoryRemoved:
		prj.parts.AddInventory(e.Part, -e.Quantity)

	case *KitAddedToProject:
		for part, quantity := range e.Parts {
			prj.parts.AddInventory(part, quantity)
		}

	case *PartsChanged:
		{
			for key, amount := range e.Removals {
				prj.parts.Remove(key, amount)
			}

			for key, amount := range e.Additions {
				prj.parts.Add(key, amount)
			}
		}

	}

}
