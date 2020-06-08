package command

import (
	"fmt"
	"mvc/app"
	"mvc/app/create"
	"os"
)

type ProjectCreateCommand struct {
	Meta
}

func (c *ProjectCreateCommand) Help() string {
	return ""
}

func (c *ProjectCreateCommand) Synopsis() string {
	return "Creates a new lego Project"
}

func (c *ProjectCreateCommand) Name() string {
	return "project create"
}

func (c *ProjectCreateCommand) Run(args []string) int {

	flags := c.FlagSet(c)

	if err := flags.Parse(args); err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	if flags.NArg() != 2 {
		c.UI.Error("This command takes two arguments, <name> and <path>")
		return 1
	}

	modelName := flags.Arg(0)
	filepath := flags.Arg(1)

	c.UI.Info(fmt.Sprintf("Creating project %s from %s", modelName, filepath))

	file, err := os.Open(filepath)
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}
	defer file.Close()

	store, err := app.NewAppStore()

	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	waiter, err := create.CreateProject(store, modelName, file)

	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	waiter()

	c.UI.Info("Done.")

	return 0
}
