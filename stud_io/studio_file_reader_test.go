package stud_io

import (
	"brickrecon/ldraw"
	"bufio"
	"os"
	"testing"

	"github.com/yeka/zip"

	"github.com/stretchr/testify/assert"
)

func TestReadingFile(t *testing.T) {

	file, err := os.Open("chaingun.io")
	assert.NoError(t, err)
	defer file.Close()

	info, err := file.Stat()
	assert.NoError(t, err)

	reader, err := zip.NewReader(file, info.Size())
	assert.NoError(t, err)

	var model *zip.File
	for _, f := range reader.File {
		if f.Name == "model.ldr" {
			model = f
			break
		}
	}

	model.SetPassword("soho0909")
	r, err := model.Open()
	assert.NoError(t, err)
	defer r.Close()

	scanner := bufio.NewScanner(r)
	bricks, err := ldraw.CreateBrickList(scanner)
	assert.NoError(t, err)
	assert.Len(t, bricks, 48)
}
