package lego

type Colour struct {
	ID      BrickLinkColour
	Aliases ColourAliases

	Name     string
	Category string
}

type BrickLinkColour int
type LDrawColour int
type BrickOwlColour int

type ColourAliases struct {
	BrickLinkID BrickLinkColour
	LDrawID     LDrawColour
	Boid        BrickOwlColour
}

type PartID string
type PartName string

type Part struct {
	ID      PartID
	Aliases PartAliases

	Name   PartName
	Colour Colour

	Quantity int
	Weight   float64
}

type PartAliases struct {
	BrickLinkID string
	ElementID   int
	LDrawID     string
	Boid        string
}
