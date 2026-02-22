package ldraw

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGeneratingMovedPartLookup(t *testing.T) {
	t.Run("direct part move", func(t *testing.T) {
		input := map[string]*partDto{
			"4073": {MovedTo: "6141"},
			"6141": {MovedTo: ""},
		}

		expected := map[string]*partDto{
			"4073": {MovedTo: "6141"},
			"6141": {MovedTo: ""},
		}

		calculateNewestMoves(input)
		require.Equal(t, expected, input)
	})

	t.Run("single chain move", func(t *testing.T) {
		input := map[string]*partDto{
			"aaa": {MovedTo: "bbb"},
			"bbb": {MovedTo: "ccc"},
			"ccc": {MovedTo: ""},
		}

		expected := map[string]*partDto{
			"aaa": {MovedTo: "ccc"},
			"bbb": {MovedTo: "ccc"},
			"ccc": {MovedTo: ""},
		}

		calculateNewestMoves(input)
		require.Equal(t, expected, input)
	})

	t.Run("multi chain move", func(t *testing.T) {
		input := map[string]*partDto{
			"aaa": {MovedTo: "bbb"},
			"bbb": {MovedTo: "ccc"},
			"ccc": {MovedTo: ""},
			"ddd": {MovedTo: "ccc"},
		}

		expected := map[string]*partDto{
			"aaa": {MovedTo: "ccc"},
			"bbb": {MovedTo: "ccc"},
			"ccc": {MovedTo: ""},
			"ddd": {MovedTo: "ccc"},
		}

		calculateNewestMoves(input)
		require.Equal(t, expected, input)
	})

	t.Run("multi chain move ordered", func(t *testing.T) {
		input := map[string]*partDto{
			"32123":  {MovedTo: "32123a"},
			"32123a": {MovedTo: ""},
			"4265c":  {MovedTo: "32123"},
		}

		expected := map[string]*partDto{
			"4265c":  {MovedTo: "32123a"},
			"32123a": {MovedTo: ""},
			"32123":  {MovedTo: "32123a"},
		}

		calculateNewestMoves(input)
		require.Equal(t, expected, input)
	})
}

func TestParsing(t *testing.T) {
	t.SkipNow()
	contents, err := os.Open("complete.zip")
	require.NoError(t, err)
	defer contents.Close()

	parts, err := ParseDatabaseArchive(t.Context(), contents)
	require.NoError(t, err)
	require.NotEmpty(t, parts)

	require.Equal(t, "6141", parts["4073"])

	require.Contains(t, parts, "4265c")
	require.Contains(t, parts, "32123")
	require.Contains(t, parts, "32123a")

	require.Equal(t, "32123a", parts["4265c"])
}

func TestDatExtraction(t *testing.T) {

	t.Run("reading part info", func(t *testing.T) {
		content, err := os.ReadFile("test_data/3626bpa8.dat")
		require.NoError(t, err)

		part := parseDat(bytes.NewReader(content))

		assert.Equal(t, "Minifig Head with Evil Skeleton Skull Pattern", part.Name)
		assert.Equal(t, "", part.MovedTo)
		assert.Equal(t, []string{"3626px115", "3626bpr0190"}, part.AlternateIds)
	})

	t.Run("reading moves", func(t *testing.T) {
		content, err := os.ReadFile("test_data/4073.dat")
		require.NoError(t, err)

		part := parseDat(bytes.NewReader(content))

		assert.Equal(t, "", part.Name)
		assert.Equal(t, "6141", part.MovedTo)
		assert.Empty(t, part.AlternateIds)
	})

}
