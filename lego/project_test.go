package lego

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddingInventory(t *testing.T) {

	partID := NewPartID("1234b")
	colourID := 5678

	project := NewProject("Test Project", []Part{
		{ID: partID, Name: "Test Part", Colour: Colour{ID: colourID}, Quantity: 5},
	})

	thePart, _ := project.FindPart(partID, colourID)

	// add to non-existing part
	assert.Error(t, project.AddInventory(NewPartID("99999"), colourID, 5))

	// add to non-existing colour
	assert.Error(t, project.AddInventory(partID, 99999, 5))

	// add negative inventory
	assert.Error(t, project.AddInventory(partID, colourID, -2))

	// add inventory
	assert.NoError(t, project.AddInventory(partID, colourID, 3))

	assert.Equal(t, 3, thePart.Inventory)
	assert.False(t, thePart.HasSpares())

	// add more inventory then quantity, gives "spares"
	assert.NoError(t, project.AddInventory(partID, colourID, 4))

	assert.Equal(t, 7, thePart.Inventory)
	assert.True(t, thePart.HasSpares())
}

func TestRemovingInventory(t *testing.T) {

	partID := NewPartID("1234b")
	colourID := 5678

	project := NewProject("Test Project", []Part{
		{ID: partID, Name: "Test Part", Colour: Colour{ID: colourID}, Quantity: 5},
	})

	project.AddInventory(partID, colourID, 4)

	thePart, _ := project.FindPart(partID, colourID)

	// remove from non-existing part
	assert.Error(t, project.RemoveInventory(NewPartID("99999"), colourID, 5))

	// remove from non-existing colour
	assert.Error(t, project.RemoveInventory(partID, 99999, 5))

	// remove negative inventory
	assert.Error(t, project.RemoveInventory(partID, colourID, -2))

	// remove inventory
	assert.NoError(t, project.RemoveInventory(partID, colourID, 2))

	assert.Equal(t, 2, thePart.Inventory)

	// remove more than inventory

	assert.NoError(t, project.RemoveInventory(partID, colourID, 17))

	assert.Equal(t, 0, thePart.Inventory)
}
