package lego

type Colour struct {
	Name     string
	Category string

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
