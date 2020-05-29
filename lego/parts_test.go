package lego

import (
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

	assert.Equal(t, "2412b", part.BrickLinkID)
	assert.Equal(t, 241226, part.ElementID)
	assert.Equal(t, "2412b", part.LDrawID)
	assert.Equal(t, "Tile, Modified 1 x 2 Grille with Bottom Groove / Lip", part.PartName)
	assert.Equal(t, 11, part.Colour.BrickLinkID)
	assert.Equal(t, 0, part.Colour.LDrawID)
	assert.Equal(t, "Black", part.Colour.Name)
	assert.Equal(t, "Solid Colors", part.Colour.Category)
	assert.Equal(t, 4, part.Quantity)
	assert.Equal(t, 0.23, part.Weight)
}
