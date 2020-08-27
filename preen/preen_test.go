package preen

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestControllerName(t *testing.T) {
	t.Parallel()

	assert.Equal(t, "single", controllerName(&singleController{}))
	assert.Equal(t, "two_name", controllerName(&TwoNameController{}))
}

type singleController struct{}

func (c singleController) Path() string    { return "" }
func (c singleController) Views() []string { return nil }

type TwoNameController struct{}

func (c TwoNameController) Path() string    { return "" }
func (c TwoNameController) Views() []string { return nil }
