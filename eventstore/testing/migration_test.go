package testing

import (
	"brickrecon/eventstore"
	"context"
	"fmt"
	"testing"

	uuid "github.com/satori/go.uuid"
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
			return event
		}

		return e
	})

	id := uuid.NewV4()
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
