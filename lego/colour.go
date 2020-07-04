package lego

type ColourName string
type HexColour string

type BrickLinkColour int
type LDrawColour int
type BrickOwlColour int

type Colour struct {
	ID      BrickLinkColour
	Aliases ColourAliases

	Name     ColourName
	Category string
	Hex      HexColour
}

type ColourAliases struct {
	BrickLinkID BrickLinkColour
	LDrawID     LDrawColour
	Boid        BrickOwlColour
}

func GetColourHex(id BrickLinkColour) HexColour {
	return HexColour(hexColours[int(id)])
}
