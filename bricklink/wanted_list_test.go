package bricklink

import (
	"brickrecon/domain"
	"brickrecon/lego"
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSerializationFormat(t *testing.T) {

	input := []*lego.InventoryPart{
		{
			Part: lego.Part{
				Id:   lego.PartId("36841"),
				Name: "testing",
			},
			ColorId:  lego.ColorId("26"),
			Quantity: 2,
		},
	}
	stock := domain.Stock{}
	domain.AddStock(stock, input[0].Id, input[0].ColorId, 1)

	expected := strings.TrimSpace(`
<INVENTORY>
  <ITEM>
    <ITEMTYPE>P</ITEMTYPE>
    <ITEMID>36841</ITEMID>
    <COLOR>11</COLOR>
    <MINQTY>2</MINQTY>
    <QTYFILLED>1</QTYFILLED>
  </ITEM>
</INVENTORY>
`)

	getPart := func(p lego.PartId) (*lego.Part, error) {
		return &lego.Part{
			Id:   p,
			Name: "testing",
		}, nil
	}

	xml, err := marshal(wantedListFromParts(input, stock))
	require.NoError(t, err)
	require.Equal(t, expected, xml)

	parts, stock, err := ParseWantedList(t.Context(), getPart, bytes.NewReader([]byte(expected)))
	require.NoError(t, err)

	require.Equal(t, input[0], parts[0])
	require.Equal(t, 1, stock["36841"]["26"], stock)
}
