package brickowl

import (
	"brickrecon/lego"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

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

	assert.Equal(t, lego.PartKey("15403|86"), part.Key)
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

func TestSanitisePartName(t *testing.T) {
	t.Parallel()

	colour := lego.Colour{
		ID:   lego.BrickLinkColour(86),
		Name: "Medium Stone Gray",
	}

	cases := map[string]string{
		"LEGO Medium Stone Gray Plate 1 x 2 with Shooter (15403b)":                          "Plate 1 x 2 with Shooter",
		"LEGO Medium Stone Gray Plate 1 x 2 with Shooter (15403b / 123132)":                 "Plate 1 x 2 with Shooter",
		"LEGO Medium Stone Gray Plate 1 x 2 with Shooter (15403)":                           "Plate 1 x 2 with Shooter",
		"LEGO Medium Stone Gray Slope 1 x 1 (31°) (50746 / 54200)":                          "Slope 1 x 1 (31°)",
		"LEGO Medium Stone Gray Dish 8 x 8 Inverted with yellow hoses and panel decoration": "Dish 8 x 8 Inverted with yellow hoses and panel decoration",
	}

	for input, expected := range cases {
		assert.Equal(t, expected, sanitisePartName(input, "15403b", colour))
	}
}

func TestBoidCsv(t *testing.T) {
	t.Parallel()

	csv := boidCsv([]lego.BrickOwlPart{lego.BrickOwlPart("123"), lego.BrickOwlPart("456")})
	assert.Equal(t, "123,456", csv)
}

func TestSanitiseKitName(t *testing.T) {
	t.Parallel()

	assert.Equal(t, "Darth Vader Transformation", string(sanitiseKitName("LEGO Darth Vader Transformation  Set 75183")))
	assert.Equal(t, "Millennium Falcon Microfighter", string(sanitiseKitName("LEGO Millennium Falcon Microfighter Set 75193")))
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
