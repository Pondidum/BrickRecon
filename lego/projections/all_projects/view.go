package all_projects

import (
	"brickrecon/lego"
	"time"

	uuid "github.com/satori/go.uuid"
)

type AllProjectsView struct {
	Names    []lego.ProjectName
	Projects map[lego.ProjectName]*ProjectView
	Kits     map[lego.KitNumber]KitView
}

type ProjectView struct {
	ID   uuid.UUID
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
	ColourID   lego.BrickLinkColour
	ColourName lego.ColourName
	ColourHex  lego.HexColour

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
	ID   lego.BrickLinkColour
	Name lego.ColourName
	Hex  lego.HexColour
}

type ProjectStatsView struct {
	TotalPartsQuantity  int
	TotalPartsInventory int
	PercentComplete     int
}
