package brickowl

import (
	"brickrecon/lego"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreatePart(t *testing.T) {
	t.Parallel()

	entry := inventoryItem{Boid: "103095-64", Quantity: 5}
	var additional lookupItem
	json.Unmarshal([]byte(partJson), &additional)

	part, err := createPart(entry, additional)

	assert.NoError(t, err)
	assert.Equal(t, lego.PartName("Plate 1 x 2 with Shooter"), part.Name)
	// assert.Equal(t, 5, part.Quantity)
	//assert.Equal(t, lego.ColorId("194"), part.ColorId)
}

func TestSanitisePartName(t *testing.T) {
	t.Parallel()

	color := "Medium Stone Gray"

	cases := map[string]string{
		"LEGO Medium Stone Gray Plate 1 x 2 with Shooter (15403b)":                          "Plate 1 x 2 with Shooter",
		"LEGO Medium Stone Gray Plate 1 x 2 with Shooter (15403b / 123132)":                 "Plate 1 x 2 with Shooter",
		"LEGO Medium Stone Gray Plate 1 x 2 with Shooter (15403)":                           "Plate 1 x 2 with Shooter",
		"LEGO Medium Stone Gray Slope 1 x 1 (31°) (50746 / 54200)":                          "Slope 1 x 1 (31°)",
		"LEGO Medium Stone Gray Dish 8 x 8 Inverted with yellow hoses and panel decoration": "Dish 8 x 8 Inverted with yellow hoses and panel decoration",
	}

	for input, expected := range cases {
		assert.Equal(t, expected, sanitisePartName(input, "15403b", color))
	}
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
