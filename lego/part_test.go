package lego

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPartIdMarshaling(t *testing.T) {

	part := NewPartID("1234b")

	bytes, err := json.Marshal(part)
	assert.NoError(t, err)
	assert.Equal(t, `"1234b"`, string(bytes))

	bytes, err = json.Marshal(container{PartID: part})
	assert.NoError(t, err)
	assert.Equal(t, `{"PartID":"1234b"}`, string(bytes))
}

func TestPartIdUnmarshaling(t *testing.T) {

	var part PartID

	err := json.Unmarshal([]byte(`"1234b"`), &part)
	assert.NoError(t, err)
	assert.Equal(t, `1234b`, part.id)

	var container container
	err = json.Unmarshal([]byte(`{"PartID":"1234b"}`), &container)
	assert.NoError(t, err)
	assert.Equal(t, "1234b", container.PartID.id)
}

func TestPartIdEquality(t *testing.T) {

	one := NewPartID("1234")
	two := NewPartID("1234")
	bad := NewPartID("9876")

	assert.True(t, one == two)
	assert.False(t, one == bad)
	assert.True(t, one != bad)
	assert.False(t, one != two)
}

func TestPartStringy(t *testing.T) {

	part := NewPartID("233b")
	assert.Equal(t, "233b", fmt.Sprintf("%s", part))
}

type container struct {
	PartID PartID
}
