package command

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTableRenderer(t *testing.T) {
	type Summary struct {
		Name string
		Age  int
	}

	items := []Summary{
		{Name: "one", Age: 31},
		{Name: "two", Age: 32},
	}

	actual := &bytes.Buffer{}
	require.NoError(t, TableRenderer(actual, items))

	expected := `Name    Age
----    ---
one     31
two     32
`
	require.Equal(t, expected, actual.String())
}
