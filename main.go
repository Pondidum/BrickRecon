package main

import (
	"fmt"
	"os"

	"mvc/command"

	"github.com/honeycombio/beeline-go"
	"github.com/honeycombio/libhoney-go"
	"github.com/honeycombio/libhoney-go/transmission"
	"github.com/mattn/go-colorable"
	"github.com/mitchellh/cli"
)

func main() {
	os.Exit(Run(os.Args[1:]))
}

func Run(args []string) int {

	beeline.Init(beeline.Config{
		// WriteKey: os.Getenv("HONEYCOMB_API_KEY"),
		STDOUT:  true,
		Dataset: "BrickRecon",
	})
	defer beeline.Close()

	err := libhoney.Init(libhoney.Config{
		APIHost:      "http://localhost",
		APIKey:       "1234",
		Transmission: &transmission.WriterSender{},
		Dataset:      "BrickRecon",
	})
	defer libhoney.Close()

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error executing CLI: %s\n", err.Error())
		return 1
	}

	ui := &cli.BasicUi{
		Reader:      os.Stdin,
		Writer:      colorable.NewColorableStdout(),
		ErrorWriter: colorable.NewColorableStderr(),
	}

	commands := command.Commands(ui)

	cli := &cli.CLI{
		Name:                       "mvc",
		Args:                       args,
		Commands:                   commands,
		Autocomplete:               true,
		AutocompleteNoDefaultFlags: false,
		HelpFunc:                   cli.BasicHelpFunc("mvc"),
		HelpWriter:                 os.Stdout,
	}

	exitCode, err := cli.Run()

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error executing CLI: %s\n", err.Error())
		return 1
	}

	return exitCode
}
