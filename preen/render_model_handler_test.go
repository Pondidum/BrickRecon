package preen

import (
	"html/template"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestControllerRegistration(t *testing.T) {

	controller := &TestImportController{}

	handler := &RenderModelHandler{
		layout: template.New("root"),
	}

	templates, err := handler.parseController("", controller)
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
