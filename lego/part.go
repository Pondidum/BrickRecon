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

	PartName string
	Colour   Colour

	Quantity int
	Weight   float64
}
