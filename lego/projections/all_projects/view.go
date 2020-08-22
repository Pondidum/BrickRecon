package all_projects

import (
	"brickrecon/lego"
	"fmt"
	"strconv"
	"strings"

	uuid "github.com/satori/go.uuid"
)

type PartKey string

func CreatePartKey(part lego.LDrawPart, colour lego.BrickLinkColour) PartKey {
	return PartKey(fmt.Sprintf("%v|%v", part, colour))
}

func ParseKey(key PartKey) (lego.LDrawPart, lego.BrickLinkColour) {
	segments := strings.Split(string(key), "|")
	val, _ := strconv.Atoi(segments[1])

	return lego.LDrawPart(segments[0]), lego.BrickLinkColour(val)
}

type AllProjectsView struct {
	Names    []lego.ProjectName
	Projects map[lego.ProjectName]*ProjectView
	Kits     map[lego.KitNumber]KitView
}

type ProjectView struct {
	ID   uuid.UUID
	Name lego.ProjectName

	Parts []*ProjectPartView
	Kits  map[lego.KitNumber]KitView

	BrickLinkXml string
}

type ProjectPartView struct {
	ID         lego.LDrawPart
	Name       lego.PartName
	ColourID   lego.BrickLinkColour
	ColourName lego.ColourName
	ColourHex  lego.HexColour

	Key PartKey

	Quantity  int
	Inventory int
}

type KitView struct {
	Number lego.KitNumber
	Name   lego.KitName

	Parts      map[PartKey]int
	TotalParts int
}
