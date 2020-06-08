package command

import (
	"testing"

	"github.com/mitchellh/cli"
	"github.com/stretchr/testify/assert"
)

func TestProjectCreation(t *testing.T) {

	meta := Meta{UI: cli.NewMockUi()}
	cmd := ProjectCreateCommand{
		Meta: meta,
	}

	exitCode := cmd.Run([]string{"testing", "../lego/test-partlist-short.csv"})

	assert.Equal(t, 0, exitCode)
}
