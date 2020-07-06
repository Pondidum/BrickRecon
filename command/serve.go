package command

import (
	"brickrecon/app"

	"net/http"
)

type ServeCommand struct {
	Meta
}

func (c *ServeCommand) Help() string {
	return ""
}

func (c *ServeCommand) Synopsis() string {
	return "Starts the server"
}

func (c *ServeCommand) Name() string {
	return "serve"
}

func (c *ServeCommand) Run(_ []string) int {

	ctx, send := c.NewPhase(c)
	defer send()

	builder := app.NewAppBuilder(ctx)
	router, err := builder.CreateWebUI()

	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	c.UI.Info("Listening on 127.0.0.1:3000")
	http.ListenAndServe("127.0.0.1:3000", router)

	return 0
}
