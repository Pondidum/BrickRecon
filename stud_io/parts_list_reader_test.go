package stud_io

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

	assert.Equal(t, lego.BrickLinkPart("2412b"), part.BrickLinkID)
	assert.Equal(t, lego.LDrawPart("2412b"), part.LDrawID)
	assert.Equal(t, lego.PartName("Tile, Modified 1 x 2 Grille with Bottom Groove / Lip"), part.Name)
	assert.Equal(t, lego.BrickLinkColour(11), part.BrickLinkColour)
	assert.Equal(t, lego.LDrawColour(0), part.LDrawColour)
	assert.Equal(t, lego.ColourName("Black"), part.ColourName)
	assert.Equal(t, "Solid Colors", part.ColourCategory)
	assert.Equal(t, 4, part.Quantity)
	assert.Equal(t, 0.23, part.Weight)
}
