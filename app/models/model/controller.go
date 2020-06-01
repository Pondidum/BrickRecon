package model

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

type ModelController struct {
	Store *app.AppStore
}

func (c ModelController) Path() string {
	return "models/{name}"
}

func (c ModelController) View() string {
	return "models/model"
}

func (c ModelController) Get(req *http.Request) interface{} {

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
