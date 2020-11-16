package lego

import (
	"encoding/json"
)

type ColoursView struct {
	ByBrickLink map[BrickLinkColour]*ColourView
	ByLDraw     map[LDrawColour]*ColourView
}

type ColourView struct {
	BrickLinkID BrickLinkColour
	LDrawID     LDrawColour
	Name        ColourName
	Hex         HexColour
	Category    string
}

func readColourView() (*ColoursView, error) {
	var view ColoursView
	if err := json.Unmarshal([]byte(colourJson), &view); err != nil {
		return nil, err
	}

	return &view, nil
}

var lookup *ColoursView

func LookupColourBricklink(bricklink int) *ColourView {
	if lookup == nil {
		v, err := readColourView()
		if err != nil {
			panic(err)
		}
		lookup = v
	}

	return lookup.ByBrickLink[BrickLinkColour(bricklink)]
}

func LookupColourLDraw(ldraw LDrawColour) (*ColourView, bool) {
	if lookup == nil {
		v, err := readColourView()
		if err != nil {
			panic(err)
		}
		lookup = v
	}

	colour, found := lookup.ByLDraw[ldraw]

	return colour, found
}

var colourJson = `
{
  "ByBrickLink": {
    "1": {
      "BrickLinkID": 1,
      "LDrawID": 15,
      "Name": "White",
      "Hex": "FFFFFF",
      "Category": "Solid Colors"
    },
    "11": {
      "BrickLinkID": 11,
      "LDrawID": 0,
      "Name": "Black",
      "Hex": "212121",
      "Category": "Solid Colors"
    },
    "115": {
      "BrickLinkID": 115,
      "LDrawID": 297,
      "Name": "Pearl Gold",
      "Hex": "e79500",
      "Category": "Pearl Colors"
    },
    "2": {
      "BrickLinkID": 2,
      "LDrawID": 19,
      "Name": "Tan",
      "Hex": "dec69c",
      "Category": "Solid Colors"
    },
    "20": {
      "BrickLinkID": 20,
      "LDrawID": 34,
      "Name": "Trans-Green",
      "Hex": "217625",
      "Category": "Transparent Colors"
    },
    "4": {
      "BrickLinkID": 4,
      "LDrawID": 25,
      "Name": "Orange",
      "Hex": "FF7E14",
      "Category": "Solid Colors"
    },
    "48": {
      "BrickLinkID": 48,
      "LDrawID": 378,
      "Name": "Sand Green",
      "Hex": "76a290",
      "Category": "Solid Colors"
    },
    "5": {
      "BrickLinkID": 5,
      "LDrawID": 4,
      "Name": "Red",
      "Hex": "b30006",
      "Category": "Solid Colors"
    },
    "59": {
      "BrickLinkID": 59,
      "LDrawID": 320,
      "Name": "Dark Red",
      "Hex": "6a0e15",
      "Category": "Solid Colors"
    },
    "66": {
      "BrickLinkID": 66,
      "LDrawID": 135,
      "Name": "Pearl Light Gray",
      "Hex": "ACB7C0",
      "Category": "Pearl Colors"
    },
    "67": {
      "BrickLinkID": 67,
      "LDrawID": 80,
      "Name": "Metallic Silver",
      "Hex": "C0C0C0",
      "Category": "Metallic Colors"
    },
    "69": {
      "BrickLinkID": 69,
      "LDrawID": 28,
      "Name": "Dark Tan",
      "Hex": "907450",
      "Category": "Solid Colors"
    },
    "7": {
      "BrickLinkID": 7,
      "LDrawID": 1,
      "Name": "Blue",
      "Hex": "0057a6",
      "Category": "Solid Colors"
    },
    "85": {
      "BrickLinkID": 85,
      "LDrawID": 72,
      "Name": "Dark Bluish Gray",
      "Hex": "595D60",
      "Category": "Solid Colors"
    },
    "86": {
      "BrickLinkID": 86,
      "LDrawID": 71,
      "Name": "Light Bluish Gray",
      "Hex": "afb5c7",
      "Category": "Solid Colors"
    },
    "88": {
      "BrickLinkID": 88,
      "LDrawID": 70,
      "Name": "Reddish Brown",
      "Hex": "89351d",
      "Category": "Solid Colors"
    },
    "95": {
      "BrickLinkID": 95,
      "LDrawID": 179,
      "Name": "Flat Silver",
      "Hex": "8D949C",
      "Category": "Pearl Colors"
    },
    "98": {
      "BrickLinkID": 98,
      "LDrawID": 57,
      "Name": "Trans-Orange",
      "Hex": "D78019",
      "Category": "Transparent Colors"
    }
  },
  "ByLDraw": {
    "0": {
      "BrickLinkID": 11,
      "LDrawID": 0,
      "Name": "Black",
      "Hex": "212121",
      "Category": "Solid Colors"
    },
    "1": {
      "BrickLinkID": 7,
      "LDrawID": 1,
      "Name": "Blue",
      "Hex": "0057a6",
      "Category": "Solid Colors"
    },
    "135": {
      "BrickLinkID": 66,
      "LDrawID": 135,
      "Name": "Pearl Light Gray",
      "Hex": "ACB7C0",
      "Category": "Pearl Colors"
    },
    "15": {
      "BrickLinkID": 1,
      "LDrawID": 15,
      "Name": "White",
      "Hex": "FFFFFF",
      "Category": "Solid Colors"
    },
    "179": {
      "BrickLinkID": 95,
      "LDrawID": 179,
      "Name": "Flat Silver",
      "Hex": "8D949C",
      "Category": "Pearl Colors"
    },
    "19": {
      "BrickLinkID": 2,
      "LDrawID": 19,
      "Name": "Tan",
      "Hex": "dec69c",
      "Category": "Solid Colors"
    },
    "25": {
      "BrickLinkID": 4,
      "LDrawID": 25,
      "Name": "Orange",
      "Hex": "FF7E14",
      "Category": "Solid Colors"
    },
    "28": {
      "BrickLinkID": 69,
      "LDrawID": 28,
      "Name": "Dark Tan",
      "Hex": "907450",
      "Category": "Solid Colors"
    },
    "297": {
      "BrickLinkID": 115,
      "LDrawID": 297,
      "Name": "Pearl Gold",
      "Hex": "e79500",
      "Category": "Pearl Colors"
    },
    "320": {
      "BrickLinkID": 59,
      "LDrawID": 320,
      "Name": "Dark Red",
      "Hex": "6a0e15",
      "Category": "Solid Colors"
    },
    "34": {
      "BrickLinkID": 20,
      "LDrawID": 34,
      "Name": "Trans-Green",
      "Hex": "217625",
      "Category": "Transparent Colors"
    },
    "378": {
      "BrickLinkID": 48,
      "LDrawID": 378,
      "Name": "Sand Green",
      "Hex": "76a290",
      "Category": "Solid Colors"
    },
    "4": {
      "BrickLinkID": 5,
      "LDrawID": 4,
      "Name": "Red",
      "Hex": "b30006",
      "Category": "Solid Colors"
    },
    "57": {
      "BrickLinkID": 98,
      "LDrawID": 57,
      "Name": "Trans-Orange",
      "Hex": "D78019",
      "Category": "Transparent Colors"
    },
    "70": {
      "BrickLinkID": 88,
      "LDrawID": 70,
      "Name": "Reddish Brown",
      "Hex": "89351d",
      "Category": "Solid Colors"
    },
    "71": {
      "BrickLinkID": 86,
      "LDrawID": 71,
      "Name": "Light Bluish Gray",
      "Hex": "afb5c7",
      "Category": "Solid Colors"
    },
    "72": {
      "BrickLinkID": 85,
      "LDrawID": 72,
      "Name": "Dark Bluish Gray",
      "Hex": "595D60",
      "Category": "Solid Colors"
    },
    "80": {
      "BrickLinkID": 67,
      "LDrawID": 80,
      "Name": "Metallic Silver",
      "Hex": "C0C0C0",
      "Category": "Metallic Colors"
    }
  }
}
`
