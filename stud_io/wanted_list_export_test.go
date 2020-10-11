package stud_io

import (
	"brickrecon/lego"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSerializationFormat(t *testing.T) {

	parts := []*lego.ProjectPart{
		{
			Part: lego.Part{
				Aliases:  lego.PartAliases{BrickLinkID: lego.BrickLinkPart("2540")},
				Colour:   lego.Colour{Aliases: lego.ColourAliases{BrickLinkID: lego.BrickLinkColour(85)}},
				Quantity: 9,
			},
			Inventory: 3,
		},
		{
			Part: lego.Part{
				Aliases:  lego.PartAliases{BrickLinkID: lego.BrickLinkPart("11477")},
				Colour:   lego.Colour{Aliases: lego.ColourAliases{BrickLinkID: lego.BrickLinkColour(59)}},
				Quantity: 2,
			},
			Inventory: 0,
		},
	}

	expected := strings.TrimSpace(`
<INVENTORY>
  <ITEM>
    <ITEMTYPE>P</ITEMTYPE>
    <ITEMID>2540</ITEMID>
    <COLOR>85</COLOR>
    <MINQTY>9</MINQTY>
    <QTYFILLED>3</QTYFILLED>
  </ITEM>
  <ITEM>
    <ITEMTYPE>P</ITEMTYPE>
    <ITEMID>11477</ITEMID>
    <COLOR>59</COLOR>
    <MINQTY>2</MINQTY>
    <QTYFILLED>0</QTYFILLED>
  </ITEM>
</INVENTORY>
`)

	xml, err := marshal(wantedListFromParts(parts))

	assert.NoError(t, err)
	assert.Equal(t, expected, xml)

}
