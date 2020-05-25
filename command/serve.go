package command

import (
	"fmt"
	"mvc/app"
	"mvc/app/create"
	"mvc/preen"

	"net/http"

	"github.com/gorilla/mux"
	"github.com/honeycombio/beeline-go/wrappers/hnynethttp"
	"github.com/mitchellh/cli"
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

	p, err := preen.NewPreen("app")

	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	r := mux.NewRouter()
	r.Use(logger(c.UI))

	p.HandleStaticAssets(r)

	r.Handle("/favicon.ico", http.NotFoundHandler())

	p.RegisterController(r, &app.AppController{})
	p.RegisterController(r, &create.CreateController{})

	c.UI.Info("Listening on 127.0.0.1:3000")
	http.ListenAndServe("127.0.0.1:3000", hnynethttp.WrapHandler(r))

	return 0
}

func logger(ui cli.Ui) mux.MiddlewareFunc {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			h.ServeHTTP(w, r)
			ui.Info(fmt.Sprintf("%s %s", r.Method, r.URL.String()))
		})
	}
}
