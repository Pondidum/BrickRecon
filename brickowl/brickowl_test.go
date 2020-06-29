package brickowl

import (
	"brickrecon/lego"
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func createApi(t *testing.T) *BrickOwlApi {
	if testing.Short() {
		t.Skip("Skipping integration tests in short mode")
	}

	return NewBrickOwlApi(os.Getenv("BRICKOWL_API_KEY"))
}

func TestFetchingBoid(t *testing.T) {
	t.Parallel()

	owl := createApi(t)
	boid, err := owl.getSetBoid("75192-1")

	assert.NoError(t, err)
	assert.Equal(t, "849212", boid)
}

func TestFetchingInventory(t *testing.T) {
	t.Parallel()

	owl := createApi(t)
	parts, err := owl.getInventory("849212")

	assert.NoError(t, err)
	assert.Len(t, parts, 716)
}

func TestBulkFetching(t *testing.T) {
	t.Parallel()

	owl := createApi(t)
	parts, err := owl.lookupParts([]string{"380995-64", "334100-64"})

	assert.NoError(t, err)
	assert.Len(t, parts, 2)
}

func TestSetLookup(t *testing.T) {
	t.Parallel()

	owl := createApi(t)
	info, err := owl.lookup("849212")

	assert.NoError(t, err)
	assert.Equal(t, "LEGO Millennium Falcon Set 75192", info.Name)
}

func TestGetInventory(t *testing.T) {
	t.Parallel()

	owl := createApi(t)
	parts, err := owl.GetParts("75193-1")

	assert.NoError(t, err)
	assert.Len(t, parts, 48)
}

func TestIdMapUnMarshal(t *testing.T) {
	t.Parallel()

	c := container{}
	data := `{ "ids": [ { "id": "4070", "type": "design_id" }, { "id": "531429-64", "type": "boid" } ] }`

	err := json.Unmarshal([]byte(data), &c)

	assert.NoError(t, err)
	assert.Contains(t, c.IDs, "design_id")
	assert.Contains(t, c.IDs, "boid")
}

type container struct {
	IDs idMap
}

func TestCreatePart(t *testing.T) {
	t.Parallel()

	colours := map[flexInt]colourItem{
		flexInt(64): {
			Name:         "Medium Stone Gray",
			ID:           "64",
			LDrawIDs:     []flexInt{flexInt(71)},
			BrickLinkIDs: []flexInt{flexInt(86)},
		},
	}

	entry := inventoryItem{Boid: "103095-64", Quantity: 5}
	var additional lookupItem
	json.Unmarshal([]byte(partJson), &additional)

	part := createPart(colours, entry, additional)

	assert.Equal(t, lego.PartID("15403"), part.ID)
	assert.Equal(t, lego.PartName("Plate 1 x 2 with Shooter"), part.Name)
	assert.Equal(t, 5, part.Quantity)
	assert.Equal(t, lego.BrickLinkColour(86), part.Colour.ID)

}

func TestCreateColour(t *testing.T) {
	t.Parallel()

	colours := map[flexInt]colourItem{
		flexInt(64): {
			Name:         "Medium Stone Gray",
			ID:           "64",
			LDrawIDs:     []flexInt{flexInt(71)},
			BrickLinkIDs: []flexInt{flexInt(86)},
		},
	}

	colour := partColour(colours, flexInt(64))

	assert.Equal(t, lego.BrickLinkColour(86), colour.ID)
	assert.Equal(t, lego.ColourName("Medium Stone Gray"), colour.Name)
	assert.Equal(t, "", colour.Category)
	assert.Equal(t, lego.BrickOwlColour(64), colour.Aliases.Boid)
	assert.Equal(t, lego.BrickLinkColour(86), colour.Aliases.BrickLinkID)
	assert.Equal(t, lego.LDrawColour(71), colour.Aliases.LDrawID)

}

func TestSanitiseName(t *testing.T) {
	t.Parallel()

	colour := lego.Colour{
		ID:   lego.BrickLinkColour(86),
		Name: "Medium Stone Gray",
	}

	cases := map[string]string{
		"LEGO Medium Stone Gray Plate 1 x 2 with Shooter (15403b)":          "Plate 1 x 2 with Shooter",
		"LEGO Medium Stone Gray Plate 1 x 2 with Shooter (15403b / 123132)": "Plate 1 x 2 with Shooter",
		"LEGO Medium Stone Gray Plate 1 x 2 with Shooter (15403)":           "Plate 1 x 2 with Shooter",
		"LEGO Medium Stone Gray Slope 1 x 1 (31°) (50746 / 54200)":          "Slope 1 x 1 (31°)",
	}

	for input, expected := range cases {
		assert.Equal(t, expected, sanitiseName(input, "15403b", colour))
	}
}

var partJson string = `
{
  "boid": "103095-64",
  "type": "Part",
  "ids": [
    { "id": "15403", "type": "design_id" },
    { "id": "103095-64", "type": "boid" },
    { "id": "6167514", "type": "item_no" },
    { "id": "6167514", "type": "item_no" }
  ],
  "name": "LEGO Medium Stone Gray Plate 1 x 2 with Shooter (15403)",
  "url": "https:\/\/www.brickowl.com\/catalog\/lego-medium-stone-gray-plate-1-x-2-with-shooter-15403",
  "permalink": "https:\/\/www.brickowl.com\/boid\/103095-64",
  "cheapest_gbp": "0.01",
  "color_name": "Medium Stone Gray",
  "color_id": "64",
  "color_hex": "afb5c7",
  "cat_name_path": "Parts \/ Plate \/ Non-Standard",
  "missing_data": "41"
}
`
