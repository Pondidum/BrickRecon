package command

import (
	"brickrecon/app"
	"brickrecon/background"
)

type ImageCacheRunCommand struct {
	Meta
}

func (c *ImageCacheRunCommand) Help() string {
	return ""
}

func (c *ImageCacheRunCommand) Synopsis() string {
	return "Fetches missing part images"
}

func (c *ImageCacheRunCommand) Name() string {
	return "imagecache run"
}

func (c *ImageCacheRunCommand) Run(args []string) int {
	ctx, send := c.NewPhase(c)
	defer send()

	store, err := app.NewAppStore(ctx)

	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	ic, err := background.NewImageCache(store.EventStore, "./app/static/img/parts", ctx)

	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	ic.Run(ctx)

	store.EventStore.SaveAggregate(ctx, ic)

	c.UI.Info("Done.")
	return 0
}
