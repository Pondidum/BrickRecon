package bricklink

import (
	"brickrecon/domain"
	"brickrecon/lego"
	"bytes"
	"context"
	"encoding/xml"
	"io"
)

type wantedList struct {
	XMLName xml.Name         `xml:"INVENTORY"`
	Items   []wantedListItem `xml:"ITEM"`
}

type wantedListItem struct {
	Type      string `xml:"ITEMTYPE"`
	ID        string `xml:"ITEMID"`
	Color     string `xml:"COLOR"`
	Quantity  int    `xml:"MINQTY"`
	Inventory int    `xml:"QTYFILLED"`
}

func AsXmlWantedList(parts []*lego.InventoryPart, stock domain.Stock) (string, error) {
	return marshal(wantedListFromParts(parts, stock))
}

func ParseWantedList(ctx context.Context, content io.Reader) ([]*lego.InventoryPart, domain.Stock, error) {
	var wantedList *wantedList

	if err := xml.NewDecoder(content).Decode(&wantedList); err != nil {
		return nil, nil, err
	}

	parts := []*lego.InventoryPart{}
	stock := domain.Stock{}

	for _, wantedItem := range wantedList.Items {
		color, err := lego.GetColorId(wantedItem.Color, "bricklink")
		if err != nil {
			return nil, nil, err
		}

		part := &lego.InventoryPart{
			Part: lego.Part{
				Id: lego.PartId(wantedItem.ID),
			},
			ColorId:  color,
			Quantity: wantedItem.Quantity,
		}

		parts = append(parts, part)

		domain.AddStock(stock, part.Id, part.ColorId, wantedItem.Inventory)
	}

	return parts, stock, nil
}

func wantedListFromParts(parts []*lego.InventoryPart, stock map[lego.PartId]map[lego.ColorId]int) wantedList {

	wanted := make([]wantedListItem, len(parts))

	getStock := func(p lego.PartId, c lego.ColorId) int {
		colorStock, found := stock[p]
		if !found {
			return 0
		}
		return colorStock[c]
	}
	for i, p := range parts {
		color, err := lego.AsColorId(p.ColorId, "bricklink")
		if err != nil {
			color = string(p.ColorId)
		}

		wanted[i] = wantedListItem{
			Type:      "P",
			ID:        string(p.Id),
			Color:     color,
			Quantity:  p.Quantity,
			Inventory: getStock(p.Id, p.ColorId),
		}
	}

	return wantedList{Items: wanted}
}

func marshal(parts wantedList) (string, error) {
	var b bytes.Buffer
	enc := xml.NewEncoder(&b)
	enc.Indent("", "  ")

	if err := enc.Encode(&parts); err != nil {
		return "", err
	}

	return b.String(), nil
}
