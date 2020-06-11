package main

import (
	"fmt"
	"os"

	"brickrecon/command"

	"github.com/honeycombio/beeline-go"
	"github.com/mattn/go-colorable"
	"github.com/mitchellh/cli"
)

func main() {
	os.Exit(Run(os.Args[1:]))
}

func Run(args []string) int {

	beeline.Init(beeline.Config{
		WriteKey: os.Getenv("HONEYCOMB_API_KEY"),
		//STDOUT:  true,
		Dataset: "BrickRecon",
	})
	defer beeline.Close()

	ui := &cli.BasicUi{
		Reader:      os.Stdin,
		Writer:      colorable.NewColorableStdout(),
		ErrorWriter: colorable.NewColorableStderr(),
	}

	commands := command.Commands(ui)

	cli := &cli.CLI{
		Name:                       "brickrecon",
		Args:                       args,
		Commands:                   commands,
		Autocomplete:               true,
		AutocompleteNoDefaultFlags: false,
		HelpFunc:                   cli.BasicHelpFunc("brickrecon"),
		HelpWriter:                 os.Stdout,
	}

	exitCode, err := cli.Run()

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error executing CLI: %s\n", err.Error())
		return 1
	}

	return exitCode
}
