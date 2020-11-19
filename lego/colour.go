package lego

import "encoding/json"

type ColourName string
type HexColour string

type BrickLinkColour int
type LDrawColour int
type BrickOwlColour int

func GetColourHex(id LDrawColour) HexColour {

	lookup := map[int]string{}

	err := json.Unmarshal([]byte(hexColours), &lookup)
	if err != nil {
		panic(err)
	}

	return HexColour(lookup[int(id)])
}
