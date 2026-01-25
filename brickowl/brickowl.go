package brickowl

import (
	"brickrecon/lego"
	"regexp"
	"strings"
)

type BrickOwlApi struct {
	api Owlette
}

func NewBrickOwlApi(key string) *BrickOwlApi {
	return &BrickOwlApi{
		api: newLowLevelApi(key),
	}
}

func (bo *BrickOwlApi) GetSetName(setNumber lego.KitNumber) (lego.KitName, error) {

	setBoid, err := bo.api.lookupSetBoid(setNumber)
	if err != nil {
		return "", err
	}

	info, err := bo.api.lookup(setBoid)
	if err != nil {
		return "", err
	}

	return sanitiseKitName(info.Name), nil
}

var rx = regexp.MustCompile(`\s+(Set\s*\d+$)`)

func sanitiseKitName(name string) lego.KitName {

	name = strings.TrimPrefix(name, "LEGO ")
	name = rx.ReplaceAllString(name, "")

	return lego.KitName(name)
}

// func (bo *BrickOwlApi) GetPart(key lego.PartKey) (*BrickOwlPart, error) {

// }

func (bo *BrickOwlApi) GetParts(setNumber lego.KitNumber) ([]*BrickOwlPart, error) {

	setBoid, err := bo.api.lookupSetBoid(setNumber)
	if err != nil {
		return nil, err
	}

	inventory, err := bo.api.getInventory(setBoid)
	if err != nil {
		return nil, err
	}

	colours, err := bo.api.listColours()
	if err != nil {
		return nil, err
	}

	chunks := split(inventory, 100)

	parts := []*BrickOwlPart{}

	for _, items := range chunks {

		partBoids := make([]lego.BrickOwlPart, len(items))

		for i, item := range items {
			partBoids[i] = item.Boid
		}

		// takes max 100 items
		partData, err := bo.api.lookupParts(partBoids)
		if err != nil {
			return nil, err
		}

		for _, item := range items {
			itemData := partData[item.Boid]

			if itemData.Type == "Part" {
				part := createPart(colours, item, itemData)
				parts = append(parts, part)
			}
		}
	}

	return parts, nil
}

func createPart(colours map[flexInt]colourItem, item inventoryItem, additional lookupItem) *BrickOwlPart {
	ldrawID, found := additional.IDs["ldraw"]

	if !found {
		ldrawID, found = additional.IDs["design_id"]
	}

	colourInfo := colours[additional.ColourID]

	name := sanitisePartName(additional.Name, ldrawID, colourInfo)
	id := lego.LDrawPart(ldrawID)
	ldColour := lego.LDrawColour(colourInfo.LDrawIDs[0])

	return &BrickOwlPart{
		Key:         lego.CreatePartKey(id, ldColour),
		Name:        lego.PartName(name),
		LDrawID:     id,
		BrickLinkID: lego.BrickLinkPart(ldrawID),
		Boid:        item.Boid,

		ColourName:      colourInfo.Name,
		ColourHex:       colourInfo.Hex,
		BrickLinkColour: lego.BrickLinkColour(colourInfo.BrickLinkIDs[0]),
		LDrawColour:     ldColour,
		ColourBoid:      lego.BrickOwlColour(additional.ColourID),

		Quantity: int(item.Quantity),
	}
}

type BrickOwlPart struct {
	Key  lego.PartKey
	Name lego.PartName

	LDrawID     lego.LDrawPart
	BrickLinkID lego.BrickLinkPart
	Boid        lego.BrickOwlPart

	ColourName      lego.ColourName
	ColourHex       lego.HexColour
	BrickLinkColour lego.BrickLinkColour
	LDrawColour     lego.LDrawColour
	ColourBoid      lego.BrickOwlColour

	Quantity int
}

func sanitisePartName(name string, id string, colour colourItem) string {

	name = strings.TrimPrefix(name, "LEGO ")
	name = strings.TrimPrefix(name, string(colour.Name))

	braceIndex := strings.LastIndex(name, "(")
	if braceIndex > 0 {
		name = name[0:strings.LastIndex(name, "(")]
	}

	name = strings.TrimSpace(name)

	return name
}

func split(buf []inventoryItem, lim int) [][]inventoryItem {
	var chunk []inventoryItem
	chunks := make([][]inventoryItem, 0, len(buf)/lim+1)
	for len(buf) >= lim {
		chunk, buf = buf[:lim], buf[lim:]
		chunks = append(chunks, chunk)
	}
	if len(buf) > 0 {
		chunks = append(chunks, buf[:])
	}
	return chunks
}
