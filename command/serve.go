package command

import (
	"fmt"
	"mvc/app"
	"mvc/app/create"
	"mvc/app/models"
	"mvc/app/models/model"
	"mvc/preen"
	"mvc/store"

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

	db := store.NewStorage()

	p, err := preen.NewPreen(preen.PreenConfig{
		ApplicationRoot: "app",
		Controllers: []preen.Controller{
			&app.AppController{DB: &db},
			&create.CreateController{DB: &db},
			&models.ModelsController{DB: &db},
			&model.ModelController{DB: &db},
		},
	})

	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	r := mux.NewRouter()
	r.Use(logger(c.UI))
	r.Use(preen.BasicAuth(preen.AuthOptions{
		Realm:    "BrickRecon",
		User:     "test",
		Password: "testing",
		Path:     "/login",
	}))

	p.Apply(r)
	r.Handle("/login", http.RedirectHandler("/", http.StatusSeeOther))

	c.UI.Info("Listening on 127.0.0.1:3000")
	http.ListenAndServe("127.0.0.1:3000", hnynethttp.WrapHandler(r))

	return 0
}

func logger(ui cli.Ui) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r)
			ui.Info(fmt.Sprintf("%s %s", r.Method, r.URL.String()))
		})
	}
}
