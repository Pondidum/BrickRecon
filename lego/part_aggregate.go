package lego

import (
	"brickrecon/eventstore"
	"strings"
)

type PartAggregate struct {
	*eventstore.Aggregator

	Key PartKey

	PartID     LDrawPart
	Name       PartName
	ColourID   LDrawColour
	ColourName ColourName

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

func BlankPart() *PartAggregate {
	part := &PartAggregate{}
	part.Aggregator = eventstore.NewAggregator(part.on)
	return part
}

func NewPart(key PartKey) *PartAggregate {
	partID, colourID := ParsePartKey(key)

	p := BlankPart()
	p.Apply(&PartCreated{Key: key, PartID: partID, ColourID: colourID})

	return p
}

func (p *PartAggregate) AddNames(partName PartName, colourName ColourName) {
	p.Apply(&PartNamesAdded{PartName: partName, ColourName: colourName})
}

func (p *PartAggregate) AddBrickOwl(boid BrickOwlPart, colourBoid BrickOwlColour) {
	p.Apply(&PartBrickOwlAdded{Part: boid, Colour: colourBoid})
}

func (p *PartAggregate) AddBrickLink(partID BrickLinkPart, colourID BrickLinkColour) {
	p.Apply(&PartBrickLinkAdded{Part: partID, Colour: colourID})
}

func (p *PartAggregate) HasImage() bool {
	return p.ImagePath != ""
}

func (p *PartAggregate) AttachImage(sourceName string, path string) {
	if strings.TrimSpace(path) == "" {
		return
	}

	p.Apply(&PartImageAdded{SourcedFrom: sourceName, Path: path})
}

func (p *PartAggregate) on(event eventstore.Event) {

	switch e := event.(type) {
	case *PartCreated:
		p.SetID(eventstore.AggregateID(string(e.Key)))
		p.Key = e.Key
		p.PartID = e.PartID
		p.ColourID = e.ColourID

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
