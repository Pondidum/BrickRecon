package lego

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddingInventory(t *testing.T) {

	partID := LDrawPart("1234b")
	colourID := BrickLinkColour(5678)
	key := CreatePartKey(partID, colourID)

	project := NewProject("Test Project", []Part{
		{ID: partID, Name: "Test Part", Colour: Colour{ID: colourID}, Quantity: 5},
	})

	thePart, _ := project.FindPart(key)

	// add to non-existing part
	assert.Error(t, project.AddInventory(CreatePartKey(LDrawPart("99999"), colourID), 5))

	// add to non-existing colour
	assert.Error(t, project.AddInventory(CreatePartKey(partID, 99999), 5))

	// add negative inventory
	assert.Error(t, project.AddInventory(key, -2))

	// add inventory
	assert.NoError(t, project.AddInventory(key, 3))

	assert.Equal(t, 3, thePart.Inventory)
	assert.False(t, thePart.HasSpares())

	// add more inventory then quantity, gives "spares"
	assert.NoError(t, project.AddInventory(key, 4))

	assert.Equal(t, 7, thePart.Inventory)
	assert.True(t, thePart.HasSpares())
}

func TestRemovingInventory(t *testing.T) {

	partID := LDrawPart("1234b")
	colourID := BrickLinkColour(5678)
	key := CreatePartKey(partID, colourID)

	project := NewProject("Test Project", []Part{
		{ID: partID, Name: "Test Part", Colour: Colour{ID: colourID}, Quantity: 5},
	})

	project.AddInventory(key, 4)

	thePart, _ := project.FindPart(key)

	// remove from non-existing part
	assert.Error(t, project.RemoveInventory(CreatePartKey(LDrawPart("99999"), colourID), 5))

	// remove from non-existing colour
	assert.Error(t, project.RemoveInventory(CreatePartKey(partID, 99999), 5))

	// remove negative inventory
	assert.Error(t, project.RemoveInventory(key, -2))

	// remove inventory
	assert.NoError(t, project.RemoveInventory(key, 2))

	assert.Equal(t, 2, thePart.Inventory)

	// remove more than inventory

	assert.NoError(t, project.RemoveInventory(key, 17))

	assert.Equal(t, 0, thePart.Inventory)
}

func TestUpdatingInventory(t *testing.T) {

	partID := LDrawPart("1234b")
	colourID := BrickLinkColour(5678)
	key := CreatePartKey(partID, colourID)

	project := NewProject("Test Project", []Part{
		{ID: partID, Name: "Test Part", Colour: Colour{ID: colourID}, Quantity: 5},
	})

	thePart, _ := project.FindPart(key)

	// no change in inventory
	project.UpdateInventory(map[PartKey]int{
		key: 5,
	})
	assert.Equal(t, 5, thePart.Inventory)

	// reduce inventory
	project.UpdateInventory(map[PartKey]int{
		key: 2,
	})
	assert.Equal(t, 2, thePart.Inventory)

	// increate inventory
	project.UpdateInventory(map[PartKey]int{
		key: 14,
	})
	assert.Equal(t, 14, thePart.Inventory)

}

func TestChangingParts(t *testing.T) {
	parts := []Part{
		createPart("123|10", 5),
		createPart("456|15", 5),
		createPart("789|10", 5),
	}

	replacementParts := []Part{
		createPart("123|10", 10),
		createPart("456|15", 5),
	}

	project := NewProject("Test Project", parts)
	project.ReplaceParts(replacementParts)

	first, _ := project.parts.FindPart(PartKey("123|10"))
	second, _ := project.parts.FindPart(PartKey("456|15"))

	assert.Len(t, project.Parts(), 2)
	assert.Equal(t, 10, first.Quantity)
	assert.Equal(t, 5, second.Quantity)
}

func createPart(key string, quantity int) Part {
	id, colour := ParsePartKey(PartKey(key))
	return Part{
		ID:       LDrawPart(id),
		Colour:   Colour{ID: BrickLinkColour(colour)},
		Quantity: quantity,
	}
}
