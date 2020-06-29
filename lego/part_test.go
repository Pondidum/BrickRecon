package lego

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPartIdMarshaling(t *testing.T) {

	part := LDrawPart("1234b")

	bytes, err := json.Marshal(part)
	assert.NoError(t, err)
	assert.Equal(t, `"1234b"`, string(bytes))

	bytes, err = json.Marshal(container{PartID: part})
	assert.NoError(t, err)
	assert.Equal(t, `{"PartID":"1234b"}`, string(bytes))
}

func TestPartIdUnmarshaling(t *testing.T) {

	var part LDrawPart

	err := json.Unmarshal([]byte(`"1234b"`), &part)
	assert.NoError(t, err)
	assert.Equal(t, LDrawPart(`1234b`), part)

	var container container
	err = json.Unmarshal([]byte(`{"PartID":"1234b"}`), &container)
	assert.NoError(t, err)
	assert.Equal(t, LDrawPart("1234b"), container.PartID)
}

func TestPartIdEquality(t *testing.T) {

	one := LDrawPart("1234")
	two := LDrawPart("1234")
	bad := LDrawPart("9876")

	assert.True(t, one == two)
	assert.False(t, one == bad)
	assert.True(t, one != bad)
	assert.False(t, one != two)
}

func TestPartStringy(t *testing.T) {

	part := LDrawPart("233b")
	assert.Equal(t, "233b", fmt.Sprintf("%s", part))
}

type container struct {
	PartID LDrawPart
}
