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

	exportedRecently bool
	lastExportMarkup string
}

func BlankProject() *Project {
	project := Project{
		parts: NewPartsList(),
	}
	project.Aggregator = eventstore.NewAggregator(project.on)

	return &project
}

func NewProject(name ProjectName, parts []*Part) *Project {

	project := BlankProject()
	project.Apply(&ProjectCreated{ID: eventstore.NewAggregateID(), Name: name})
	project.Apply(&ProjectPartsAdded{Parts: parts})

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

	partID, colourID := ParsePartKey(part)
	prj.Apply(&ProjectInventoryAdded{PartID: partID, ColourID: colourID, Quantity: quantity})

	return nil
}

func (prj *Project) RemoveInventory(part PartKey, quantity int) error {

	if quantity <= 0 {
		return errors.New("Quantity must be greater than 0")
	}

	if _, found := prj.parts.FindPart(part); !found {
		return fmt.Errorf("No part with id %s found", part)
	}

	partID, colourID := ParsePartKey(part)
	prj.Apply(&ProjectInventoryRemoved{PartID: partID, ColourID: colourID, Quantity: quantity})

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
			prj.Apply(&ProjectInventoryRemoved{PartID: partID, ColourID: colourID, Quantity: diff * -1})
		}

		if diff > 0 {
			prj.Apply(&ProjectInventoryAdded{PartID: partID, ColourID: colourID, Quantity: diff})
		}
	}

	return nil
}

func (prj *Project) AddKitContents(number KitNumber, name KitName, parts []PartQuantity) {

	if len(parts) == 0 {
		return
	}

	prj.Apply((&KitAddedToProject{KitNumber: number, KitName: name, Parts: parts}))
}

func (prj *Project) ReplaceParts(parts []*Part) map[PartKey]int {
	other := NewPartsList()
	for _, part := range parts {
		other.Add(part)
	}

	changes := prj.parts.Diff(other)

	if len(changes) == 0 {
		return changes
	}

	event := &PartsChanged{
		Additions: []*Part{},
		Removals:  map[PartKey]int{},
	}

	for key, change := range changes {
		if change < 0 {
			event.Removals[key] = change * -1
		}

		if change > 0 {
			part, _ := other.FindPart(key)
			part.Quantity = change
			event.Additions = append(event.Additions, part.Part)
		}
	}

	prj.Apply(event)

	return changes
}

func (prj *Project) Parts() []*ProjectPart {
	return prj.parts.All()
}

func (prj *Project) Diff(parts []*Part) map[PartKey]int {
	other := NewPartsList()
	for _, part := range parts {
		other.Add(part)
	}

	return prj.parts.Diff(other)
}

func (prj *Project) ExportWantedList(exporter Exporter) (string, error) {

	if prj.exportedRecently {
		return prj.lastExportMarkup, nil
	}

	wanted := []*ProjectPart{}

	for _, part := range prj.parts.parts {
		if part.IsFulfilled() == false {
			wanted = append(wanted, part)
		}
	}

	markup, err := exporter.Export(wanted)
	if err != nil {
		return "", err
	}

	prj.Apply(&WantedListExported{
		Type:   exporter.GetExporterType(),
		Markup: markup,
	})

	return markup, nil
}

type Exporter interface {
	GetExporterType() string

	Export(parts []*ProjectPart) (string, error)
}

func (prj *Project) on(event eventstore.Event) {

	prj.exportedRecently = false

	switch e := event.(type) {

	case *ProjectCreated:
		prj.SetID(e.ID)
		prj.Name = e.Name

	case *ProjectPartsAdded:
		for _, p := range e.Parts {
			prj.parts.Add(p)
		}

	case *ProjectInventoryAdded:
		prj.parts.AddInventory(CreatePartKey(e.PartID, e.ColourID), e.Quantity)

	case *ProjectInventoryRemoved:
		prj.parts.AddInventory(CreatePartKey(e.PartID, e.ColourID), -e.Quantity)

	case *KitAddedToProject:
		for _, pq := range e.Parts {
			prj.parts.AddInventory(CreatePartKey(pq.PartID, pq.ColourID), pq.Quantity)
		}

	case *PartsChanged:
		{
			for key, amount := range e.Removals {
				prj.parts.Remove(key, amount)
			}

			for _, part := range e.Additions {
				prj.parts.Add(part)
			}
		}

	case *WantedListExported:
		prj.exportedRecently = true
		prj.lastExportMarkup = e.Markup

	}

}
