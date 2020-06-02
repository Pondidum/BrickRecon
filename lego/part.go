package lego

type Colour struct {
	Name     string
	Category string

	BrickLinkID int
	LDrawID     int
}

type Part struct {
	BrickLinkID string
	ElementID   int
	LDrawID     string

	Name   string
	Colour Colour

	Quantity int
	Weight   float64
}
