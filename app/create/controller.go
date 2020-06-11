package create

import (
	"brickrecon/app"
	"brickrecon/preen"
	"net/http"
)

type CreateController struct {
	Store *app.AppStore
}

func (c CreateController) Path() string {
	return "create"
}

func (c CreateController) AuthRequired() bool {
	return true
}

func (c CreateController) Get(req *http.Request) interface{} {
	return c.Store.SiteModel()
}

func (c CreateController) Post(req *http.Request) interface{} {
	file, _, err := req.FormFile("modelFile")
	modelName := req.FormValue("modelName")

	if err != nil {
		return preen.ComposeModels(c.Store.SiteModel(), preen.ErrorModel(err))
	}

	defer file.Close()

	_, err = CreateProject(req.Context(), c.Store, modelName, file)

	if err != nil {
		return preen.ComposeModels(c.Store.SiteModel(), preen.ErrorModel(err))
	}

	return preen.Redirect{URL: "/project/" + modelName}
}
