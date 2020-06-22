package preen

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTemplateNaming(t *testing.T) {
	t.Parallel()

	testCases := map[string]string{
		"_shared/link.html":        "link",
		"_shared/image/circle.svg": "image/circle",
	}

	for path, name := range testCases {
		assert.Equal(t, name, templateName("_shared", path))
	}

	controllerCases := map[string]string{
		"create_index.html":    "create",
		"create_quantity.html": "create/quantity",
	}

	for path, name := range controllerCases {
		assert.Equal(t, name, templateName("create", path))
	}
}
