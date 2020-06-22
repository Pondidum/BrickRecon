package app

import (
	"brickrecon/preen"
	"net/http"
)

type CreateController struct {
	Store *AppStore
}

func (c CreateController) Views() []string {
	return []string{
		"create_index.html",
		"create_form.html",
	}
}

func (c CreateController) Path() string {
	return "create"
}

func (c CreateController) AuthRequired() bool {
	return true
}

func (c CreateController) Get(req *http.Request) interface{} {
	return c.Store.SiteModel(req.Context())
}

func (c CreateController) Post(req *http.Request) interface{} {
	ctx := req.Context()
	file, _, err := req.FormFile("modelFile")
	modelName := req.FormValue("modelName")

	if err != nil {
		return preen.ComposeModels(c.Store.SiteModel(ctx), preen.ErrorModel(err))
	}

	defer file.Close()

	_, err = CreateProject(ctx, c.Store, modelName, file)

	if err != nil {
		return preen.ComposeModels(c.Store.SiteModel(ctx), preen.ErrorModel(err))
	}

	return preen.Redirect{URL: "/project/" + modelName}
}
