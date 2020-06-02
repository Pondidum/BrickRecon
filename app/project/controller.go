package project

import (
	"mvc/app"
	"mvc/lego"
	"mvc/preen"
	"net/http"

	"github.com/gorilla/mux"
)

type ProjectModel struct {
	Project *lego.ProjectView
}

type ProjectController struct {
	Store *app.AppStore
}

func (c ProjectController) Path() string {
	return "project/{name}"
}

func (c ProjectController) View() string {
	return "project"
}

func (c ProjectController) Get(req *http.Request) interface{} {

	vars := mux.Vars(req)

	siteModel := c.Store.SiteModel()
	selected, _ := c.Store.Project(vars["name"])

	return preen.ComposeModels(
		siteModel,
		ProjectModel{
			Project: selected,
		},
	)
}
