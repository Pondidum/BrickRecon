package adapters

import (
	"brickrecon/lego"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPartReading(t *testing.T) {
	reader, err := os.Open("test-partlist.csv")
	assert.NoError(t, err)

	partList, err := ReadPartsList(reader)

	assert.NoError(t, err)
	assert.Len(t, partList, 128)

	part := partList[8]

	assert.Equal(t, lego.NewPartID("2412b"), part.ID)
	assert.Equal(t, "2412b", part.Aliases.BrickLinkID)
	assert.Equal(t, 241226, part.Aliases.ElementID)
	assert.Equal(t, "2412b", part.Aliases.LDrawID)
	assert.Equal(t, "Tile, Modified 1 x 2 Grille with Bottom Groove / Lip", part.Name)
	assert.Equal(t, 11, part.Colour.ID)
	assert.Equal(t, 11, part.Colour.Aliases.BrickLinkID)
	assert.Equal(t, 0, part.Colour.Aliases.LDrawID)
	assert.Equal(t, "Black", part.Colour.Name)
	assert.Equal(t, "Solid Colors", part.Colour.Category)
	assert.Equal(t, 4, part.Quantity)
	assert.Equal(t, 0.23, part.Weight)
}
