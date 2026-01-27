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

	assert.Equal(t, lego.PartId("2412b"), part.Id)
	assert.Equal(t, lego.PartName("Tile, Modified 1 x 2 Grille with Bottom Groove / Lip"), part.Name)
	assert.Equal(t, lego.ColorId("26"), part.ColourId)
	assert.Equal(t, 4, part.Quantity)
}
