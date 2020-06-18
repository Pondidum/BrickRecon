package command

import (
	"brickrecon/app"
)

type EventStoreViewsRebuildCommand struct {
	Meta
}

func (c *EventStoreViewsRebuildCommand) Help() string {
	return ""
}

func (c *EventStoreViewsRebuildCommand) Synopsis() string {
	return "Rebuilds the EventStore's views"
}

func (c *EventStoreViewsRebuildCommand) Name() string {
	return "eventstore view rebuild"
}

func (c *EventStoreViewsRebuildCommand) Run(args []string) int {
	ctx, send := c.NewPhase(c)
	defer send()

	builder := app.NewAppStoreBuilder(ctx)

	backend, err := builder.CreateBackend()

	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	store := builder.CreateEventStore(backend)

	if err := backend.DestroyViews(); err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	if err := store.RunProjections(ctx); err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	c.UI.Info("Done.")
	return 0
}
