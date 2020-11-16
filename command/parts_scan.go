package command

import (
	"brickrecon/app"
	"brickrecon/eventstore"
	"brickrecon/lego"
	"fmt"

	"github.com/honeycombio/beeline-go"
)

type PartsScanCommand struct {
	Meta
}

func (c *PartsScanCommand) Help() string {
	return ""
}

func (c *PartsScanCommand) Synopsis() string {
	return "Scans Projects for missing parts, and adds them"
}

func (c *PartsScanCommand) Name() string {
	return "parts scan"
}

func (c *PartsScanCommand) Run(args []string) int {

	ctx, send := c.NewPhase(c)
	defer send()

	flags := c.FlagSet(c)

	if err := flags.Parse(args); err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	if flags.NArg() != 0 {
		c.UI.Error("This command takes no arguments")
		return 1
	}

	store, err := app.NewAppBuilder(ctx).CreateAppStore()
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	builder, err := app.NewPartBuilder(ctx, store.EventStore)
	if err != nil {
		beeline.AddField(ctx, "builder_error", err)
		c.UI.Error(err.Error())
		return 1
	}

	store.EventStore.RunStatelessProjection(ctx, eventstore.NewStatelessProjection("scan_parts", func(e eventstore.Event) {

		switch event := e.(type) {

		case *lego.ProjectPartsAdded:
			for _, part := range event.Parts {
				err := builder.StorePart(ctx, part)

				if err != nil {
					c.UI.Info(fmt.Sprintf("Failed processing %s", part.Key))
					beeline.AddField(ctx, string(part.Key)+"_error", err)
				} else {
					c.UI.Info(fmt.Sprintf("Processed %s", part.Key))
				}
			}
		}
	}))

	beeline.AddField(ctx, "complete", true)
	c.UI.Info("Done.")

	return 0
}
