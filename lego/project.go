package lego

import (
	"brickrecon/eventstore"
	"context"
	"errors"
	"fmt"

	uuid "github.com/satori/go.uuid"
)

type Project struct {
	*eventstore.Aggregator

	Name  string
	parts *ProjectPartList
}

func NewProject(name string, parts []Part) *Project {

	project := Project{
		parts: &ProjectPartList{},
	}
	project.Aggregator = eventstore.NewAggregator(project.on)

	project.Apply(&ProjectCreated{ID: uuid.NewV4(), Name: name})
	project.Apply(&ProjectPartsAdded{Parts: parts})

	return &project
}

func (prj *Project) FindPart(partID string, colourID int) (*ProjectPart, bool) {
	return prj.parts.FindPart(partID, colourID)
}

func (prj *Project) AddInventory(partID string, colourID int, quantity int) error {

	if quantity <= 0 {
		return errors.New("Quantity must be greater than 0")
	}

	if _, found := prj.parts.FindPart(partID, colourID); !found {
		return fmt.Errorf("No part with id %s and colour %v found", partID, colourID)
	}

	prj.Apply(&ProjectInventoryAdded{PartID: partID, ColourID: colourID, Quantity: quantity})

	return nil
}

func (prj *Project) RemoveInventory(partID string, colourID int, quantity int) error {

	if quantity <= 0 {
		return errors.New("Quantity must be greater than 0")
	}

	if _, found := prj.parts.FindPart(partID, colourID); !found {
		return fmt.Errorf("No part with id %s and colour %v found", partID, colourID)
	}

	prj.Apply(&ProjectInventoryRemoved{PartID: partID, ColourID: colourID, Quantity: quantity})

	return nil
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

	}

}

type ProjectCreated struct {
	eventstore.EventMeta

	ID   uuid.UUID
	Name string
}

type ProjectPartsAdded struct {
	eventstore.EventMeta

	Parts []Part
}

type ProjectInventoryAdded struct {
	eventstore.EventMeta

	PartID   string
	ColourID int
	Quantity int
}

type ProjectInventoryRemoved struct {
	eventstore.EventMeta

	PartID   string
	ColourID int
	Quantity int
}

func ProjectEvents(ctx context.Context, register func(context.Context, eventstore.Initialiser)) {
	register(ctx, func() interface{} { return &ProjectCreated{} })
	register(ctx, func() interface{} { return &ProjectPartsAdded{} })
	register(ctx, func() interface{} { return &ProjectInventoryAdded{} })
	register(ctx, func() interface{} { return &ProjectInventoryRemoved{} })
}
