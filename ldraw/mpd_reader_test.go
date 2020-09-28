package ldraw

import (
	"bufio"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadingFile(t *testing.T) {

	file, err := os.Open("kazutortype.mpd")
	assert.NoError(t, err)
	reader := bufio.NewScanner(file)

	models, err := parseFile(reader)

	assert.NoError(t, err)
	assert.Equal(t, "kazutortype", models["kazutortype"].Name)
	assert.Equal(t, 0, models["kazutortype"].Index)

	assert.Len(t, models, 45)

	assert.Len(t, models["kazutortype"].Models, 5)
	assert.Equal(t, 16, models["kazutortype"].Models[0].PrimaryColour)
	assert.Equal(t, "shoulder right", models["kazutortype"].Models[0].Name)

	assert.Equal(t, 1, models["shoulder right"].Index)
	assert.Len(t, models["shoulder right"].Parts, 21)
	assert.Len(t, models["shoulder right"].Models, 4)

	bricks, err := collectBricks(models)

	assert.NoError(t, err)
	assert.Len(t, bricks, 233)
}
