package lego

import (
	"brickrecon/eventstore"

	uuid "github.com/satori/go.uuid"
)

type KitNumber string
type KitName string

type Kit struct {
	*eventstore.Aggregator

	KitNumber KitNumber
	Name      string

	parts []*Part
}

func BlankKit() *Kit {
	kit := Kit{}
	kit.Aggregator = eventstore.NewAggregator(kit.on)

	return &kit
}

func ImportKit(number KitNumber, name KitName, parts []*Part) *Kit {
	kit := BlankKit()
	kit.Apply(&KitCreated{
		ID:        eventstore.AggregateID(uuid.NewV4().String()),
		KitNumber: number,
		KitName:   name,
		Parts:     parts,
	})

	return kit
}

func (kit *Kit) on(event eventstore.Event) {

	switch e := event.(type) {

	case *KitCreated:
		kit.SetID(e.ID)
		kit.KitNumber = e.KitNumber
		kit.parts = e.Parts
	}
}

type KitCreated struct {
	eventstore.EventMeta

	ID        eventstore.AggregateID
	KitNumber KitNumber
	KitName   KitName
	Parts     []*Part
}

var KitEvents = []eventstore.Initialiser{
	func() interface{} { return &KitCreated{} },
}
