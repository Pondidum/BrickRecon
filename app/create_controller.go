package app

import (
	"brickrecon/lego"
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

func (c CreateController) Get(pc *preen.PreenContext, req *http.Request) interface{} {
	return nil
}

func (c CreateController) Post(pc *preen.PreenContext, req *http.Request) interface{} {
	ctx := req.Context()
	file, _, err := req.FormFile("modelFile")
	modelName := lego.ProjectName(req.FormValue("modelName"))

	if err != nil {
		return preen.ErrorModel(err)
	}

	defer file.Close()

	_, err = CreateProject(ctx, c.Store, modelName, file)

	if err != nil {
		return preen.ErrorModel(err)
	}

	return preen.ControllerRedirect("project", "name", string(modelName))
}
