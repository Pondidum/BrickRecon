package ldraw

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGeneratingMovedPartLookup(t *testing.T) {
	t.Run("direct part move", func(t *testing.T) {
		input := map[string]string{
			"4073": "6141",
			"6141": "",
		}

		expected := map[string]string{
			"4073": "6141",
			"6141": "",
		}

		result := buildParts(input)
		require.Equal(t, expected, result)
	})

	t.Run("single chain move", func(t *testing.T) {
		input := map[string]string{
			"aaa": "bbb",
			"bbb": "ccc",
			"ccc": "",
		}

		expected := map[string]string{
			"aaa": "ccc",
			"bbb": "ccc",
			"ccc": "",
		}

		result := buildParts(input)
		require.Equal(t, expected, result)
	})

	t.Run("multi chain move", func(t *testing.T) {
		input := map[string]string{
			"aaa": "bbb",
			"bbb": "ccc",
			"ccc": "",
			"ddd": "ccc",
		}

		expected := map[string]string{
			"aaa": "ccc",
			"bbb": "ccc",
			"ccc": "",
			"ddd": "ccc",
		}

		result := buildParts(input)
		require.Equal(t, expected, result)
	})

	t.Run("multi chain move ordered", func(t *testing.T) {
		input := map[string]string{
			"32123":  "32123a",
			"32123a": "",
			"4265c":  "32123",
		}

		expected := map[string]string{
			"4265c":  "32123a",
			"32123a": "",
			"32123":  "32123a",
		}

		result := buildParts(input)
		require.Equal(t, expected, result)
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
