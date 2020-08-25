package stud_io

import (
	"brickrecon/lego"
	"bytes"
	"encoding/xml"
)

type wantedList struct {
	XMLName xml.Name         `xml:"INVENTORY"`
	Items   []wantedListItem `xml:"ITEM"`
}

type wantedListItem struct {
	Type      string               `xml:"ITEMTYPE"`
	ID        lego.BrickLinkPart   `xml:"ITEMID"`
	Color     lego.BrickLinkColour `xml:"COLOR"`
	Quantity  int                  `xml:"MINQTY"`
	Inventory int                  `xml:"QTYFILLED"`
}

type WantedListXmlExporter struct{}

func (e *WantedListXmlExporter) GetExporterType() string {
	return "BrickLinkWantedListXml"
}

func (e *WantedListXmlExporter) Export(parts []*lego.ProjectPart) (string, error) {
	return marshal(wantedListFromParts(parts))
}

func wantedListFromParts(parts []*lego.ProjectPart) wantedList {

	wanted := make([]wantedListItem, len(parts))

	for i, p := range parts {
		wanted[i] = wantedListItem{
			Type:      "P",
			ID:        p.Aliases.BrickLinkID,
			Color:     p.Colour.Aliases.BrickLinkID,
			Quantity:  p.Quantity,
			Inventory: p.NeededQuantity(),
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
