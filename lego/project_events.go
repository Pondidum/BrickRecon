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

	Parts map[PartKey]int
}

type ProjectInventoryAdded struct {
	eventstore.EventMeta

	Part     PartKey
	Quantity int
}

type ProjectInventoryRemoved struct {
	eventstore.EventMeta

	Part     PartKey
	Quantity int
}

type KitAddedToProject struct {
	eventstore.EventMeta

	KitNumber KitNumber
	KitName   KitName
	Parts     map[PartKey]int
}

type PartsChanged struct {
	eventstore.EventMeta

	Additions map[PartKey]int
	Removals  map[PartKey]int
}

var ProjectEvents = []eventstore.Initialiser{
	func() interface{} { return &ProjectCreated{} },
	func() interface{} { return &ProjectPartsAdded{} },
	func() interface{} { return &ProjectInventoryAdded{} },
	func() interface{} { return &ProjectInventoryRemoved{} },
	func() interface{} { return &KitAddedToProject{} },
	func() interface{} { return &PartsChanged{} },
}
