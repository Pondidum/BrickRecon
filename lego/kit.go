package lego

import (
	"brickrecon/eventstore"
	"context"

	uuid "github.com/satori/go.uuid"
)

type Kit struct {
	*eventstore.Aggregator

	KitNumber string
	Name      string

	// parts []Part
}

func BlankKit() *Kit {
	kit := Kit{}
	kit.Aggregator = eventstore.NewAggregator(kit.on)

	return &kit
}

func NewKit(kitNumber string) *Kit {
	kit := BlankKit()
	kit.Apply(&KitCreated{ID: uuid.NewV4(), KitNumber: kitNumber})
	return kit
}

func (kit *Kit) on(event eventstore.Event) {

	switch e := event.(type) {

	case *KitCreated:
		kit.SetID(e.ID)
		kit.KitNumber = e.KitNumber
	}
}

type KitCreated struct {
	eventstore.EventMeta

	ID        uuid.UUID
	KitNumber string
}

func KitEvents(ctx context.Context, register func(context.Context, eventstore.Initialiser)) {
	register(ctx, func() interface{} { return &KitCreated{} })
}
