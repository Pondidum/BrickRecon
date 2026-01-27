package brickowl

import (
	"brickrecon/lego"
	"fmt"
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

func (bo *BrickOwlApi) GetSetName(setNumber lego.SetId) (lego.SetName, error) {

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

func sanitiseKitName(name string) lego.SetName {

	name = strings.TrimPrefix(name, "LEGO ")
	name = rx.ReplaceAllString(name, "")

	return lego.SetName(name)
}

func (bo *BrickOwlApi) GetParts(setNumber lego.SetId) ([]*lego.InventoryPart, error) {

	setBoid, err := bo.api.lookupSetBoid(setNumber)
	if err != nil {
		return nil, err
	}

	inventory, err := bo.api.getInventory(setBoid)
	if err != nil {
		return nil, err
	}

	chunks := split(inventory, 100)

	parts := []*lego.InventoryPart{}

	for _, items := range chunks {

		partBoids := make([]Boid, len(items))

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
				part, err := createPart(item, itemData)
				if err != nil {
					return nil, err
				}

				color, err := lego.GetColorId(fmt.Sprint(itemData.ColorID), "brickowl")
				if err != nil {
					return nil, err
				}
				parts = append(parts, &lego.InventoryPart{
					Part:    part,
					ColorId: color,

					Quantity: int(item.Quantity),
				})
			}
		}
	}

	return parts, nil
}

func createPart(item inventoryItem, additional lookupItem) (lego.Part, error) {
	ldrawID, found := additional.IDs["ldraw"]

	if !found {
		ldrawID, found = additional.IDs["design_id"]
	}

	name := sanitisePartName(additional.Name, ldrawID, additional.ColorName)
	id := lego.PartId(ldrawID)

	return lego.Part{
		Name: lego.PartName(name),
		Id:   id,
		Sources: []lego.Source{
			{SourceName: "brickowl", PartId: string(item.Boid)},
		},
	}, nil
}

func sanitisePartName(name string, id string, colorName string) string {

	name = strings.TrimPrefix(name, "LEGO ")
	name = strings.TrimPrefix(name, colorName)

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
