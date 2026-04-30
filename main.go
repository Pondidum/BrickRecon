package main

import (
	"brickrecon/command"
	"brickrecon/command/version"
	"fmt"
	"os"

	"github.com/hashicorp/cli"
)

func main() {

	defs := []command.CommandDefinition{
		version.NewVersionCommand(),
		command.NewProjectNewCommand(),
		command.NewProjectListCommand(),
		command.NewProjectViewCommand(),
		command.NewProjectRmCommand(),
		command.NewProjectPartsImportCommand(),
		command.NewProjectSetAddCommand(),
		command.NewProjectSetFindCommand(),
		command.NewDatabaseSyncCommand(),
	}

	commands := make(map[string]cli.CommandFactory, len(defs))
	for _, def := range defs {
		commands[def.Name()] = command.NewCommand(def)
	}

	cli := &cli.CLI{
		Name:                       "kirjasto",
		Args:                       os.Args[1:],
		Commands:                   commands,
		Autocomplete:               true,
		AutocompleteNoDefaultFlags: false,
	}

	exitCode, err := cli.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error executing CLI: %s\n", err.Error())
	}

	os.Exit(exitCode)
}
