package lego

type LDrawPart string
type PartName string

type Part struct {
	Key     PartKey
	Aliases PartAliases

	Name   PartName
	Colour Colour

	Quantity int
	Weight   float64
}

type PartQuantity struct {
	Part     PartKey
	PartID   LDrawPart
	ColourID BrickLinkColour
	Quantity int
}

type BrickLinkPart string
type BrickOwlPart string

type PartAliases struct {
	BrickLinkID BrickLinkPart
	LDrawID     LDrawPart
	Boid        BrickOwlPart
}
