package command

import (
	"brickrecon/version"

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
		"parts scan": func() (cli.Command, error) {
			return &PartsScanCommand{Meta: meta}, nil
		},
		"project create": func() (cli.Command, error) {
			return &ProjectCreateCommand{Meta: meta}, nil
		},
		"project list": func() (cli.Command, error) {
			return &ProjectListCommand{Meta: meta}, nil
		},
		"project replace": func() (cli.Command, error) {
			return &ProjectReplaceCommand{Meta: meta}, nil
		},
		"kit import": func() (cli.Command, error) {
			return &KitImportCommand{Meta: meta}, nil
		},
		"eventstore view rebuild": func() (cli.Command, error) {
			return &EventStoreViewsRebuildCommand{Meta: meta}, nil
		},
	}

	return all
}
