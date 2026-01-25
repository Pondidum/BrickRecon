package main

import (
	"fmt"
	"os"

	"github.com/hashicorp/cli"
)

func main() {
	os.Exit(Run(os.Args[1:]))
}

func Run(args []string) int {

	// ui := &cli.ColoredUi{
	// 	WarnColor:  cli.UiColorYellow,
	// 	ErrorColor: cli.UiColorRed,
	// 	Ui: &cli.BasicUi{
	// 		Reader:      os.Stdin,
	// 		Writer:      colorable.NewColorableStdout(),
	// 		ErrorWriter: colorable.NewColorableStderr(),
	// 	},
	// }

	commands := map[string]cli.CommandFactory{}

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
