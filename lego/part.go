package lego

type PartName string
type LDrawPart string
type BrickLinkPart string
type BrickOwlPart string

type Part struct {
	PartID     LDrawPart
	Name       PartName
	ColourID   LDrawColour
	ColourName ColourName
	ColourHex  HexColour

	ImagePath string

	BrickOwl  BrickOwl
	BrickLink BrickLink
}

type BrickOwl struct {
	PartBoid   BrickOwlPart
	ColourBoid BrickOwlColour
}

type BrickLink struct {
	PartNumber BrickLinkPart
	Colour     BrickLinkColour
}

func NewPart(partId LDrawPart, name PartName) *Part {
	return &Part{PartID: partId, Name: name}
}
func (p *Part) HasImage() bool {
	return p.ImagePath != ""
}
