package lego

import (
	"brickrecon/eventstore"
	"strings"
)

type PartAggregate struct {
	*eventstore.Aggregator

	Key PartKey

	Number     LDrawPart
	Name       PartName
	Colour     LDrawColour
	ColourName ColourName

	ImagePath string

	BrickOwl  BrickOwl
	BrickLink BrickLink
}

type BrickOwl struct {
	ID BrickOwlID

	PartNumber *BrickOwlPart
	Colour     *BrickOwlColour
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

func NewPartFromLDraw(key PartKey, id LDrawPart, name PartName, colour LDrawColour, colourName ColourName, colourCategory string) *PartAggregate {
	p := BlankPart()
	p.Apply(&PartCreated{Key: key})
	p.Apply(&PartLDrawAdded{PartID: id, Name: name, Colour: colour, ColourName: colourName, ColourCategory: colourCategory})

	return p
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

	case *PartLDrawAdded:
		p.Number = e.PartID
		p.Name = e.Name
		p.Colour = e.Colour
		p.ColourName = e.ColourName

	case *PartImageAdded:
		p.ImagePath = e.Path
	}

}

type PartCreated struct {
	eventstore.EventMeta

	Key PartKey
}

type PartLDrawAdded struct {
	eventstore.EventMeta

	PartID         LDrawPart
	Name           PartName
	Colour         LDrawColour
	ColourName     ColourName
	ColourCategory string
}

type PartImageAdded struct {
	eventstore.EventMeta

	Path        string
	SourcedFrom string
}

var PartEvents = []eventstore.Initialiser{
	func() interface{} { return &PartCreated{} },
	func() interface{} { return &PartLDrawAdded{} },
	func() interface{} { return &PartImageAdded{} },
}
