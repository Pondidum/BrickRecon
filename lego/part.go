package lego

import (
	"brickrecon/eventstore"
	"strings"
)

type PartName string
type LDrawPart string
type BrickLinkPart string
type BrickOwlPart string

type Part struct {
	*eventstore.Aggregator

	Key PartKey

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

func BlankPart() *Part {
	part := &Part{}
	part.Aggregator = eventstore.NewAggregator(part.on)
	return part
}

func NewPart(key PartKey) *Part {
	partID, colourID := ParsePartKey(key)

	p := BlankPart()
	p.Apply(&PartCreated{Key: key, PartID: partID, ColourID: colourID})

	return p
}

func (p *Part) AddNames(partName PartName, colourName ColourName) {
	p.Apply(&PartNamesAdded{PartName: partName, ColourName: colourName})
}

func (p *Part) AddBrickOwl(boid BrickOwlPart, colourBoid BrickOwlColour) {
	p.Apply(&PartBrickOwlAdded{Part: boid, Colour: colourBoid})
}

func (p *Part) AddBrickLink(partID BrickLinkPart, colourID BrickLinkColour) {
	p.Apply(&PartBrickLinkAdded{Part: partID, Colour: colourID})
}

func (p *Part) HasImage() bool {
	return p.ImagePath != ""
}

func (p *Part) AttachImage(sourceName string, path string) {
	if strings.TrimSpace(path) == "" {
		return
	}

	p.Apply(&PartImageAdded{SourcedFrom: sourceName, Path: path})
}

func (p *Part) on(event eventstore.Event) {

	switch e := event.(type) {
	case *PartCreated:
		p.SetID(eventstore.AggregateID(string(e.Key)))
		p.Key = e.Key
		p.PartID = e.PartID
		p.ColourID = e.ColourID
		p.ColourHex = GetColourHex(p.ColourID)

	case *PartNamesAdded:
		p.Name = e.PartName
		p.ColourName = e.ColourName

	case *PartBrickOwlAdded:
		p.BrickOwl = BrickOwl{
			PartBoid:   e.Part,
			ColourBoid: e.Colour,
		}

	case *PartBrickLinkAdded:
		p.BrickLink = BrickLink{
			PartNumber: e.Part,
			Colour:     e.Colour,
		}

	case *PartImageAdded:
		p.ImagePath = e.Path
	}

}

type PartCreated struct {
	eventstore.EventMeta

	Key      PartKey
	PartID   LDrawPart
	ColourID LDrawColour
}

type PartNamesAdded struct {
	eventstore.EventMeta

	PartName   PartName
	ColourName ColourName
}

type PartBrickOwlAdded struct {
	eventstore.EventMeta

	Part   BrickOwlPart
	Colour BrickOwlColour
}

type PartBrickLinkAdded struct {
	eventstore.EventMeta

	Part   BrickLinkPart
	Colour BrickLinkColour
}

type PartImageAdded struct {
	eventstore.EventMeta

	Path        string
	SourcedFrom string
}

var PartEvents = []eventstore.Initialiser{
	func() interface{} { return &PartCreated{} },
	func() interface{} { return &PartNamesAdded{} },
	func() interface{} { return &PartBrickOwlAdded{} },
	func() interface{} { return &PartBrickLinkAdded{} },
	func() interface{} { return &PartImageAdded{} },
}
