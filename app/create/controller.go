package create

import (
	"mvc/app"
	"mvc/preen"
	"net/http"
)

type CreateController struct{}

func (c CreateController) Path() string {
	return "create"
}

func (c CreateController) Get(req *http.Request) interface{} {
	return app.SiteModel{Models: []string{"one", "two", "three"}}
}

func (c CreateController) Post(req *http.Request) interface{} {
	// _, handler, _ := req.FormFile("modelFile")
	// fileName := req.FormValue("modelName")

	// c.UI.Info(fmt.Sprintf("Create Model: %s, %s (%v)", fileName, handler.Filename, handler.Size))

	// model := app.SiteModel{Models: []string{"one", "two", "three", fileName}}

	return preen.Redirect{URL: "/"}
}
