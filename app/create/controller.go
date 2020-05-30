package create

import (
	"mvc/app"
	"mvc/lego"
	"mvc/preen"
	"mvc/store"
	"net/http"
)

type CreateModel struct {
	ErrorMessage string
}

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
	return app.SiteModel{AllModels: c.DB.GetModelNames()}
}

func (c CreateController) Post(req *http.Request) interface{} {
	file, _, err := req.FormFile("modelFile")
	modelName := req.FormValue("modelName")

	if err != nil {
		return preen.ComposeModels(app.SiteModel{AllModels: c.DB.GetModelNames()}, CreateModel{
			ErrorMessage: err.Error(),
		})
	}

	defer file.Close()

	parts, err := lego.ReadPartsList(file)

	if err != nil {
		return preen.ComposeModels(app.SiteModel{AllModels: c.DB.GetModelNames()}, CreateModel{
			ErrorMessage: err.Error(),
		})
	}

	legoModel := lego.NewProject(modelName, parts)
	c.DB.AddModel(legoModel)

	return preen.Redirect{URL: "/models/" + modelName}
}
