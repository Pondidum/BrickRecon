package all_projects

import (
	"brickrecon/eventstore"
	"brickrecon/lego"
	"time"
)

type AllProjectsView struct {
	Names    []lego.ProjectName
	Projects map[lego.ProjectName]*ProjectView
	Kits     map[lego.KitNumber]KitView
}

type ProjectView struct {
	ID   eventstore.AggregateID
	Name lego.ProjectName

	Parts   []*ProjectPartView
	Kits    map[lego.KitNumber]KitView
	Colours []*ColourView

	Stats *ProjectStatsView

	BrickLinkXml string

	Events []*EventDescription
}

type EventDescription struct {
	Timestamp   time.Time
	Type        string
	Description string
	Additional  map[string]interface{}
}

func (e *EventDescription) With(key string, value interface{}) *EventDescription {
	e.Additional[key] = value
	return e
}

type ProjectPartView struct {
	Key        lego.PartKey
	ID         lego.LDrawPart
	Name       lego.PartName
	ColourID   lego.LDrawColour
	ColourName lego.ColourName
	ColourHex  lego.HexColour

	ImagePath string

	Quantity  int
	Inventory int
}

type KitView struct {
	Number lego.KitNumber
	Name   lego.KitName

	Parts      map[lego.PartKey]int
	TotalParts int
}

type ColourView struct {
	ID   lego.LDrawColour
	Name lego.ColourName
	Hex  lego.HexColour
}

type ProjectStatsView struct {
	TotalPartsQuantity  int
	TotalPartsInventory int
	PercentComplete     int
}
