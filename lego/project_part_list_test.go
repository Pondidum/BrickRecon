package lego

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPartListAdding(t *testing.T) {

	model := NewPartsList([]Part{})
	assert.Len(t, model.parts, 0)

	// add a part
	model.Add(Part{
		ID:       PartID("1234"),
		Colour:   Colour{ID: 1, Name: "Black"},
		Quantity: 1,
	})
	assert.Len(t, model.parts, 1)

	// duplicate part should increase quantity
	model.Add(Part{
		ID:       PartID("1234"),
		Colour:   Colour{ID: 1, Name: "Black"},
		Quantity: 17,
	})
	assert.Len(t, model.parts, 1)
	assert.Equal(t, 18, model.parts[0].Quantity)

	// duplicate part with differnt colour
	model.Add(Part{
		ID:       PartID("1234"),
		Colour:   Colour{ID: 2, Name: "Red"},
		Quantity: 1,
	})
	assert.Len(t, model.parts, 2)
	assert.Equal(t, 18, model.parts[0].Quantity)
	assert.Equal(t, 1, model.parts[1].Quantity)
}
