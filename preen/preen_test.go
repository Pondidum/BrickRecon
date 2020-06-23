package preen

import (
	"html/template"
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

func TestControllerRegistration(t *testing.T) {

	controller := &TestImportController{}

	parent := template.New("root")

	templates, err := parseController("", parent, controller)
	assert.NoError(t, err)
	assert.Contains(t, templates, "test/import")
	assert.Contains(t, templates, "test/import/form")
	assert.Contains(t, templates, "test/import/dir")
	assert.Contains(t, templates, "test/import/dir/test")
}

type TestImportController struct{}

func (c TestImportController) Path() string { return "test/import" }
func (c TestImportController) Views() []string {
	return []string{
		"test_import_index.html",
		"test_import_form.html",
		"dir/index.html",
		"dir/test.html",
	}
}
