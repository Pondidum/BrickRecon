package brickowl

import (
	"brickrecon/lego"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func createApi(t *testing.T) Owlette {
	if testing.Short() {
		t.Skip("Skipping integration tests in short mode")
	}

	return newLowLevelApi(os.Getenv("BRICKOWL_API_KEY"))
}

func TestFetchingBoid(t *testing.T) {
	t.Parallel()

	owl := createApi(t)
	boid, err := owl.lookupSetBoid("75192-1")

	assert.NoError(t, err)
	assert.Equal(t, "849212", boid)
}

func TestFetchingInventory(t *testing.T) {
	t.Parallel()

	owl := createApi(t)
	parts, err := owl.getInventory("849212")

	assert.NoError(t, err)
	assert.Len(t, parts, 716)
}

func TestBulkFetching(t *testing.T) {
	t.Parallel()

	owl := createApi(t)
	parts, err := owl.lookupParts([]lego.BrickOwlPart{"380995-64", "334100-64"})

	assert.NoError(t, err)
	assert.Len(t, parts, 2)
}

func TestSetLookup(t *testing.T) {
	t.Parallel()

	owl := createApi(t)
	info, err := owl.lookup("849212")

	assert.NoError(t, err)
	assert.Equal(t, "LEGO Millennium Falcon Set 75192", info.Name)
}

func TestBoidCsv(t *testing.T) {
	t.Parallel()

	csv := boidCsv([]lego.BrickOwlPart{lego.BrickOwlPart("123"), lego.BrickOwlPart("456")})
	assert.Equal(t, "123,456", csv)
}
