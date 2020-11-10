package app

import (
	"brickrecon/eventstore"
	"brickrecon/lego"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEventMigration(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	t.Run("ProjectInventoryAdded", func(t *testing.T) {
		added := &lego.ProjectInventoryAdded{
			EventMeta: eventstore.EventMeta{},
			PartID:    lego.LDrawPart("123"),
			ColourID:  lego.LDrawColour(67),
		}

		NewAppBuilder(ctx).upgradeEvent(ctx, added)
		assert.Equal(t, 1, added.EventVersion)
		assert.Equal(t, lego.PartKey("123|80"), added.Part)
	})

	t.Run("ProjectInventoryRemoved", func(t *testing.T) {
		added := &lego.ProjectInventoryRemoved{
			EventMeta: eventstore.EventMeta{},
			PartID:    lego.LDrawPart("123"),
			ColourID:  lego.LDrawColour(67),
		}

		NewAppBuilder(ctx).upgradeEvent(ctx, added)
		assert.Equal(t, 1, added.EventVersion)
		assert.Equal(t, lego.PartKey("123|80"), added.Part)
	})

}
