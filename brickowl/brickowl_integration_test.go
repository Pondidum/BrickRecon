package brickowl

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetInventory(t *testing.T) {
	t.Parallel()

	if testing.Short() {
		t.Skip("Skipping integration tests in short mode")
	}

	owl := NewBrickOwlApi(os.Getenv("BRICKOWL_API_KEY"))

	parts, err := owl.GetParts("75193-1")

	assert.NoError(t, err)
	assert.Len(t, parts, 48)
}
