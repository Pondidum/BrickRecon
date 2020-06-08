package command

import (
	"mvc/version"

	"github.com/mitchellh/cli"
)

func Commands(ui cli.Ui) map[string]cli.CommandFactory {

	meta := Meta{UI: ui}

	all := map[string]cli.CommandFactory{
		"version": func() (cli.Command, error) {
			return &VersionCommand{Meta: meta, Version: version.GetVersion()}, nil
		},
		"serve": func() (cli.Command, error) {
			return &ServeCommand{Meta: meta}, nil
		},
		"project create": func() (cli.Command, error) {
			return &ProjectCreateCommand{Meta: meta}, nil
		},
	}

	return all
}
