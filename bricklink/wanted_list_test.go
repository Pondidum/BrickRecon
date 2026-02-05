package bricklink

import (
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
				Id: lego.PartId("36841"),
			},
			ColorId:  lego.ColorId("26"),
			Quantity: 2,
		},
	}

	expected := strings.TrimSpace(`
<INVENTORY>
  <ITEM>
    <ITEMTYPE>P</ITEMTYPE>
    <ITEMID>36841</ITEMID>
    <COLOR>11</COLOR>
    <MINQTY>2</MINQTY>
    <QTYFILLED>0</QTYFILLED>
  </ITEM>
</INVENTORY>
`)

	xml, err := marshal(wantedListFromParts(input, nil))
	require.NoError(t, err)
	require.Equal(t, expected, xml)

	parts, stock, err := ParseWantedList(t.Context(), bytes.NewReader([]byte(expected)))
	require.NoError(t, err)

	require.Equal(t, input[0], parts[0])
	require.Equal(t, 2, stock["36841"]["26"])
}
