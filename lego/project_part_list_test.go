package lego

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPartListAdding(t *testing.T) {

	partID := LDrawPart("1234")
	black := LDrawColour(1)
	red := LDrawColour(2)

	model := NewPartsList()
	assert.Len(t, model.parts, 0)

	// add a part
	model.Add(&Part{
		Key: CreatePartKey(partID, black),

		Quantity: 1,
	})
	assert.Len(t, model.parts, 1)

	// duplicate part should increase quantity
	model.Add(&Part{
		Key: CreatePartKey(partID, black),

		Quantity: 17,
	})
	assert.Len(t, model.parts, 1)
	assert.Equal(t, 18, model.parts[CreatePartKey(partID, black)].Quantity)

	// duplicate part with differnt colour
	model.Add(&Part{
		Key: CreatePartKey(partID, red),

		Quantity: 1,
	})
	assert.Len(t, model.parts, 2)
	assert.Equal(t, 18, model.parts[CreatePartKey(partID, black)].Quantity)
	assert.Equal(t, 1, model.parts[CreatePartKey(partID, red)].Quantity)
}

func TestDiffingPartLists(t *testing.T) {
	t.Parallel()

	t.Run("Identical lists", func(t *testing.T) {
		start := partList(map[PartKey]int{
			"123|10": 5,
		})

		updated := partList(map[PartKey]int{
			"123|10": 5,
		})

		assert.Empty(t, start.Diff(updated))
		assert.Empty(t, updated.Diff(start))
	})

	t.Run("Part Quantity Increase", func(t *testing.T) {
		start := partList(map[PartKey]int{
			"123|10": 5,
		})

		updated := partList(map[PartKey]int{
			"123|10": 8,
		})

		assert.Equal(t, map[PartKey]int{
			PartKey("123|10"): 3,
		}, start.Diff(updated))
	})

	t.Run("Part Quantity Decrease", func(t *testing.T) {
		start := partList(map[PartKey]int{
			"123|10": 5,
		})

		updated := partList(map[PartKey]int{
			"123|10": 1,
		})

		assert.Equal(t, map[PartKey]int{
			PartKey("123|10"): -4,
		}, start.Diff(updated))
	})

	t.Run("Remove a part", func(t *testing.T) {
		start := partList(map[PartKey]int{
			"123|10": 5,
		})

		updated := partList(map[PartKey]int{})

		assert.Equal(t, map[PartKey]int{
			PartKey("123|10"): -5,
		}, start.Diff(updated))
	})

	t.Run("Add a part", func(t *testing.T) {
		start := partList(map[PartKey]int{
			"123|10": 5,
		})

		updated := partList(map[PartKey]int{
			"123|10": 5,
			"456|14": 2,
		})

		assert.Equal(t, map[PartKey]int{
			PartKey("456|14"): 2,
		}, start.Diff(updated))
	})

}

func partList(parts map[PartKey]int) *ProjectPartList {

	list := NewPartsList()

	for key, quantity := range parts {
		list.Add(&Part{
			Key:      key,
			Quantity: quantity,
		})
	}

	return list
}
