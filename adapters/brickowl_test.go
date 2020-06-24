package adapters

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func createApi(t *testing.T) *BrickOwlApi {
	if testing.Short() {
		t.Skip("Skipping integration tests in short mode")
	}

	return NewBrickOwlApi(os.Getenv("BRICKOWL_API_KEY"))
}

func TestFetchingBoid(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	owl := createApi(t)
	boid, err := owl.getSetBoid("75192-1")

	assert.NoError(t, err)
	assert.Equal(t, "849212", boid)
}

func TestFetchingInventory(t *testing.T) {

	owl := createApi(t)
	parts, err := owl.getInventory("849212")

	assert.NoError(t, err)
	assert.Len(t, parts, 716)
}

func TestBulkFetching(t *testing.T) {

	owl := createApi(t)
	parts, err := owl.lookupParts([]string{"380995-64", "334100-64"})

	assert.NoError(t, err)
	assert.Len(t, parts, 2)
}

func TestSetLookup(t *testing.T) {
	owl := createApi(t)
	info, err := owl.lookup("849212")

	assert.NoError(t, err)
	assert.Equal(t, "LEGO Millennium Falcon Set 75192", info.Name)
}

func TestGetInventory(t *testing.T) {

	owl := createApi(t)
	parts, err := owl.GetParts("75193-1")

	assert.NoError(t, err)
	assert.Len(t, parts, 48)
}

func TestIdMapUnMarshal(t *testing.T) {
	c := container{}
	data := `{ "ids": [ { "id": "4070", "type": "design_id" }, { "id": "531429-64", "type": "boid" } ] }`

	err := json.Unmarshal([]byte(data), &c)

	assert.NoError(t, err)
	assert.Contains(t, c.IDs, "design_id")
	assert.Contains(t, c.IDs, "boid")
}

type container struct {
	IDs idMap
}
