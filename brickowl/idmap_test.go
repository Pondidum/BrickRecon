package brickowl

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIdMapUnMarshal(t *testing.T) {
	t.Parallel()

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
