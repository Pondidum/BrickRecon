package lego

import (
	"brickrecon/eventstore"
)

type ProjectCreated struct {
	eventstore.EventMeta

	ID   eventstore.AggregateID
	Name ProjectName
}

type ProjectPartsAdded struct {
	eventstore.EventMeta

	Parts []*Part
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

type WantedListExported struct {
	eventstore.EventMeta

	Type   string
	Markup string
}

type PartsChanged struct {
	eventstore.EventMeta

	Additions []*Part
	Removals  map[PartKey]int
}

var ProjectEvents = []eventstore.Initialiser{
	func() interface{} { return &ProjectCreated{} },
	func() interface{} { return &ProjectPartsAdded{} },
	func() interface{} { return &ProjectInventoryAdded{} },
	func() interface{} { return &ProjectInventoryRemoved{} },
	func() interface{} { return &KitAddedToProject{} },
	func() interface{} { return &WantedListExported{} },
	func() interface{} { return &PartsChanged{} },
}
