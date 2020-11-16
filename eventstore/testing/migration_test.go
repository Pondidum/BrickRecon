package testing

import (
	"brickrecon/eventstore"
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEventMigration(t *testing.T) {

	be, cleanup := createBackend()
	defer cleanup()

	es := eventstore.NewEventStore(be)
	es.RegisterEvent(context.Background(), func() interface{} { return &MigrationTestEvent{} })
	es.RegisterEventMiddleware(context.Background(), func(c context.Context, e eventstore.Event) eventstore.Event {

		switch event := e.(type) {
		case *MigrationTestEvent:
			event.EventVersion = 1
			event.Key = fmt.Sprintf("%v:%v", event.Part, event.Colour)
		}

		return e
	})

	id := eventstore.NewAggregateID()
	a := eventstore.NewAggregator(func(e eventstore.Event) {})
	a.SetID(id)
	a.Apply(&MigrationTestEvent{Part: 123, Colour: 456})

	assert.NoError(t, es.SaveAggregate(context.Background(), a))

	loaded := BlankLoggingAggregate()
	assert.NoError(t, es.LoadAggregate(context.Background(), id, loaded.Aggregator))

	assert.Len(t, loaded.Applied, 1)
	e := loaded.Applied[0]
	m := e.(*MigrationTestEvent)

	assert.Equal(t, 123, m.Part)
	assert.Equal(t, 456, m.Colour)
	assert.Equal(t, 1, m.EventVersion)
	assert.Equal(t, "123:456", m.Key)
}

func TestEventMigrationTyped(t *testing.T) {

	be, cleanup := createBackend()
	defer cleanup()

	es := eventstore.NewEventStore(be)
	es.RegisterEvent(context.Background(), func() interface{} { return &InputEvent{} })
	es.RegisterEvent(context.Background(), func() interface{} { return &OutputEvent{} })

	es.RegisterEventMiddleware(context.Background(), func(c context.Context, e eventstore.Event) eventstore.Event {

		switch event := e.(type) {
		case *InputEvent:
			return &OutputEvent{
				EventMeta: event.EventMeta,
				Name:      fmt.Sprintf("%v|%v", event.PartID, event.ColourID),
			}

		}

		return e
	})

	id := eventstore.NewAggregateID()
	a := eventstore.NewAggregator(func(e eventstore.Event) {})
	a.SetID(id)
	a.Apply(&InputEvent{PartID: 123, ColourID: 456})

	assert.NoError(t, es.SaveAggregate(context.Background(), a))

	loaded := BlankLoggingAggregate()
	assert.NoError(t, es.LoadAggregate(context.Background(), id, loaded.Aggregator))

	assert.Len(t, loaded.Applied, 1)
	e := loaded.Applied[0]
	m := e.(*OutputEvent)

	assert.Equal(t, "123|456", m.Name)
}

type InputEvent struct {
	eventstore.EventMeta

	PartID   int
	ColourID int
}

type OutputEvent struct {
	eventstore.EventMeta

	Name string
}
