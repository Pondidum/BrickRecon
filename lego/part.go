package lego

type Colour struct {
	ID      int
	Aliases ColourAliases

	Name     string
	Category string
}

type ColourAliases struct {
	BrickLinkID int
	LDrawID     int
}

type Part struct {
	ID      string
	Aliases PartAliases

	Name   string
	Colour Colour

	Quantity int
	Weight   float64
}

type PartAliases struct {
	BrickLinkID string
	ElementID   int
	LDrawID     string
}
