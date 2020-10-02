package lego

import (
	"brickrecon/eventstore"
	"errors"
	"fmt"

	uuid "github.com/satori/go.uuid"
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

func NewProject(name ProjectName, parts []Part) *Project {

	project := BlankProject()
	project.Apply(&ProjectCreated{ID: uuid.NewV4(), Name: name})
	project.Apply(&ProjectPartsAdded{Parts: parts})

	return project
}

func (prj *Project) FindPart(partID LDrawPart, colourID BrickLinkColour) (*ProjectPart, bool) {
	return prj.parts.FindPart(partID, colourID)
}

func (prj *Project) AddInventory(partID LDrawPart, colourID BrickLinkColour, quantity int) error {

	if quantity <= 0 {
		return errors.New("Quantity must be greater than 0")
	}

	if _, found := prj.parts.FindPart(partID, colourID); !found {
		return fmt.Errorf("No part with id %s and colour %v found", partID, colourID)
	}

	prj.Apply(&ProjectInventoryAdded{PartID: partID, ColourID: colourID, Quantity: quantity})

	return nil
}

func (prj *Project) RemoveInventory(partID LDrawPart, colourID BrickLinkColour, quantity int) error {

	if quantity <= 0 {
		return errors.New("Quantity must be greater than 0")
	}

	if _, found := prj.parts.FindPart(partID, colourID); !found {
		return fmt.Errorf("No part with id %s and colour %v found", partID, colourID)
	}

	prj.Apply(&ProjectInventoryRemoved{PartID: partID, ColourID: colourID, Quantity: quantity})

	return nil
}

func (prj *Project) UpdateInventory(inventoryState map[PartKey]int) error {

	for key, quantity := range inventoryState {

		partID, colourID := ParsePartKey(key)

		current, found := prj.parts.FindPartByKey(key)
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
		prj.parts.AddInventory(e.PartID, e.ColourID, e.Quantity)

	case *ProjectInventoryRemoved:
		prj.parts.AddInventory(e.PartID, e.ColourID, -e.Quantity)

	case *KitAddedToProject:
		for _, pq := range e.Parts {
			prj.parts.AddInventory(pq.PartID, pq.ColourID, pq.Quantity)
		}

	case *WantedListExported:
		prj.exportedRecently = true
		prj.lastExportMarkup = e.Markup

	}

}
