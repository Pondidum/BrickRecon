package command

import (
	"brickrecon/actions"
	"brickrecon/app"
	"brickrecon/lego"
	"fmt"

	"github.com/honeycombio/beeline-go"
)

type KitImportCommand struct {
	Meta
}

func (c *KitImportCommand) Help() string {
	return ""
}

func (c *KitImportCommand) Synopsis() string {
	return "Imports a Lego set"
}

func (c *KitImportCommand) Name() string {
	return "kit import"
}

func (c *KitImportCommand) Run(args []string) int {

	ctx, send := c.NewPhase(c)
	defer send()

	flags := c.FlagSet(c)

	if err := flags.Parse(args); err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	if flags.NArg() != 1 {
		c.UI.Error("This command takes one argument, <kitnumber>")
		return 1
	}

	kitNumber := lego.KitNumber(flags.Arg(0))

	c.UI.Info(fmt.Sprintf("Importing kit number %s", kitNumber))

	store, err := app.NewAppBuilder(ctx).CreateAppStore()
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	waiter, err := actions.ImportKit(ctx, store.EventStore, kitNumber)

	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	waiter()

	beeline.AddField(ctx, "complete", true)

	c.UI.Info("Done.")

	return 0
}
