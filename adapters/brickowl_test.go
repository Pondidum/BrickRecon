package adapters

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFetchingBoid(t *testing.T) {

	owl := NewBrickOwlApi("46ee1ad3d0cf66d4d41be4b92c2923c99c84d85ced698b553a955a896e851124")
	boid, err := owl.getSetBoid("75192-1")

	assert.NoError(t, err)
	assert.Equal(t, "849212", boid)
}

func TestFetchingInventory(t *testing.T) {

	owl := NewBrickOwlApi("46ee1ad3d0cf66d4d41be4b92c2923c99c84d85ced698b553a955a896e851124")
	parts, err := owl.getInventory("849212")

	assert.NoError(t, err)
	assert.Len(t, parts, 716)
}

func TestBulkFetching(t *testing.T) {

	owl := NewBrickOwlApi("46ee1ad3d0cf66d4d41be4b92c2923c99c84d85ced698b553a955a896e851124")
	parts, err := owl.lookupParts([]string{"380995-64", "334100-64"})

	assert.NoError(t, err)
	assert.Len(t, parts, 2)
}

func TestSetLookup(t *testing.T) {
	owl := NewBrickOwlApi("46ee1ad3d0cf66d4d41be4b92c2923c99c84d85ced698b553a955a896e851124")
	info, err := owl.lookup("849212")

	assert.NoError(t, err)
	assert.Equal(t, "LEGO Millennium Falcon Set 75192", info.Name)
}

func TestGetInventory(t *testing.T) {

	owl := NewBrickOwlApi("46ee1ad3d0cf66d4d41be4b92c2923c99c84d85ced698b553a955a896e851124")
	parts, err := owl.GetParts("75193-1")

	assert.NoError(t, err)
	assert.Len(t, parts, 48)
}
