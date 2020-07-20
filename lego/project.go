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
}

func BlankProject() *Project {
	project := Project{
		parts: &ProjectPartList{},
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

func (prj *Project) AddKitContents(number KitNumber, name KitName, parts []PartQuantity) {

	if len(parts) == 0 {
		return
	}

	prj.Apply((&KitAddedToProject{KitNumber: number, KitName: name, Parts: parts}))
}

func (prj *Project) on(event eventstore.Event) {

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

	}

}

type ProjectCreated struct {
	eventstore.EventMeta

	ID   uuid.UUID
	Name ProjectName
}

type ProjectPartsAdded struct {
	eventstore.EventMeta

	Parts []Part
}

type ProjectInventoryAdded struct {
	eventstore.EventMeta

	PartID   LDrawPart
	ColourID BrickLinkColour
	Quantity int
}

type ProjectInventoryRemoved struct {
	eventstore.EventMeta

	PartID   LDrawPart
	ColourID BrickLinkColour
	Quantity int
}

type KitAddedToProject struct {
	eventstore.EventMeta

	KitNumber KitNumber
	KitName   KitName
	Parts     []PartQuantity
}

var ProjectEvents = []eventstore.Initialiser{
	func() interface{} { return &ProjectCreated{} },
	func() interface{} { return &ProjectPartsAdded{} },
	func() interface{} { return &ProjectInventoryAdded{} },
	func() interface{} { return &ProjectInventoryRemoved{} },
	func() interface{} { return &KitAddedToProject{} },
}
