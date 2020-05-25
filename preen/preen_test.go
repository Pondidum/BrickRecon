package preen

import (
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testCases map[string]string = map[string]string{
	"_shared/link.html":        "link",
	"_shared/image/circle.svg": "image/circle",
	"menu/index.html":          "menu",
	"menu/list.html":           "menu/list",
	"menu/child/index.html":    "menu/child",
	"menu/child/item.html":     "menu/child/item",
}

func TestTemplateNaming(t *testing.T) {
	t.Parallel()

	for path, name := range testCases {
		assert.Equal(t, name, templateName(path))
	}
}

func TestAreaLoading(t *testing.T) {

	temp, _ := ioutil.TempDir(".", "app")

	defer func() {
		os.RemoveAll(temp)
	}()

	for p, v := range testCases {
		os.MkdirAll(path.Join(temp, path.Dir(p)), os.ModePerm)
		ioutil.WriteFile(path.Join(temp, p), []byte(v), 0644)
	}

	ioutil.WriteFile(path.Join(temp, "index.html"), []byte("root"), 0644)

	p, err := NewPreen(PreenConfig{ApplicationRoot: temp})

	assert.NoError(t, err)

	for _, v := range testCases {
		assert.Contains(t, p.templates, v)
	}
}
