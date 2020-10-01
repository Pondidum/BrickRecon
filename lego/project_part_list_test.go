package lego

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPartListAdding(t *testing.T) {

	partID := LDrawPart("1234")
	black := BrickLinkColour(1)
	red := BrickLinkColour(2)

	model := NewPartsList()
	assert.Len(t, model.parts, 0)

	// add a part
	model.Add(Part{
		ID:       partID,
		Colour:   Colour{ID: black, Name: "Black"},
		Quantity: 1,
	})
	assert.Len(t, model.parts, 1)

	// duplicate part should increase quantity
	model.Add(Part{
		ID:       partID,
		Colour:   Colour{ID: black, Name: "Black"},
		Quantity: 17,
	})
	assert.Len(t, model.parts, 1)
	assert.Equal(t, 18, model.parts[CreatePartKey(partID, black)].Quantity)

	// duplicate part with differnt colour
	model.Add(Part{
		ID:       partID,
		Colour:   Colour{ID: red, Name: "Red"},
		Quantity: 1,
	})
	assert.Len(t, model.parts, 2)
	assert.Equal(t, 18, model.parts[CreatePartKey(partID, black)].Quantity)
	assert.Equal(t, 1, model.parts[CreatePartKey(partID, red)].Quantity)
}
