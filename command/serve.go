package command

import (
	"mvc/preen"

	"net/http"

	"github.com/gorilla/mux"
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

	r.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		p.View(w, req, DashboardModel{Models: []string{"one", "two", "three"}})
	})

	r.HandleFunc("/{area}", func(w http.ResponseWriter, req *http.Request) {
		p.View(w, req, DashboardModel{Models: []string{"one", "two", "three"}})
	})

	c.UI.Info("Listening on 127.0.0.1:3000")
	http.ListenAndServe("127.0.0.1:3000", r)

	return 0
}

func logger(ui cli.Ui) mux.MiddlewareFunc {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ui.Info(r.URL.String())
			h.ServeHTTP(w, r)
		})
	}
}

type DashboardModel struct {
	Models []string
}
