package fs

import (
	"brickrecon/eventstore"
	"brickrecon/lego"
	"context"
	"testing"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func TestReadingLongLines(t *testing.T) {
	t.Parallel()

	registry := eventstore.NewRegistry()
	registry.Register(context.Background(), func() interface{} { return &lego.KitCreated{} })

	reader, err := NewAggregateEventReader(
		context.Background(),
		registry,
		DirectoryPath("testcases"),
		uuid.Must(uuid.FromString("8347bc62-c8f4-492c-ad3a-96c33fc52b2a")),
	)

	assert.NoError(t, err)
	defer reader.Close()

	events := []eventstore.Event{}
	for reader.Read() {
		e, err := reader.Event()
		assert.NoError(t, err)

		events = append(events, e)
	}

	assert.Len(t, events, 1)
}
