package create

import (
	"mvc/app"
	"mvc/preen"
	"mvc/store"
	"net/http"
)

type CreateController struct {
	DB *store.Storage
}

func (c CreateController) Path() string {
	return "create"
}

func (c CreateController) AuthRequired() bool {
	return true
}

func (c CreateController) Get(req *http.Request) interface{} {
	return app.SiteModel{Models: c.DB.GetModels()}
}

func (c CreateController) Post(req *http.Request) interface{} {
	// _, handler, _ := req.FormFile("modelFile")
	modelName := req.FormValue("modelName")

	// c.UI.Info(fmt.Sprintf("Create Model: %s, %s (%v)", fileName, handler.Filename, handler.Size))

	// model := app.SiteModel{Models: []string{"one", "two", "three", fileName}}

	c.DB.AddModel(modelName)

	return preen.Redirect{URL: "/models/" + modelName}
}
