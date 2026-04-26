package command

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

type testSummary struct {
	Name string
	Age  int
}

func TestTableRenderer(t *testing.T) {

	items := []testSummary{
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

func TestJsonRenderer(t *testing.T) {

	items := []testSummary{
		{Name: "one", Age: 31},
		{Name: "two", Age: 32},
	}

	actual := &bytes.Buffer{}
	require.NoError(t, JsonRenderer(actual, items))

	expected := `[
  {
    "Name": "one",
    "Age": 31
  },
  {
    "Name": "two",
    "Age": 32
  }
]
`
	require.Equal(t, expected, actual.String())
}
