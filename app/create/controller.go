package create

import (
	"mvc/app"
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
	fileName := req.FormValue("modelName")

	// c.UI.Info(fmt.Sprintf("Create Model: %s, %s (%v)", fileName, handler.Filename, handler.Size))

	return app.SiteModel{Models: []string{"one", "two", "three", fileName}}
	// p.View(w, req, DashboardModel{Models: []string{"one", "two", "three"}})
}
