package project

import (
	"brickrecon/app"
	"brickrecon/lego"
	"brickrecon/preen"
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

	siteModel := c.Store.SiteModel(req.Context())
	selected, _ := c.Store.ReadProject(req.Context(), vars["name"])

	return preen.ComposeModels(
		siteModel,
		ProjectModel{
			Project: selected,
		},
	)
}

func (c ProjectController) Post(req *http.Request) interface{} {
	vars := mux.Vars(req)

	siteModel := c.Store.SiteModel(req.Context())
	selected, _ := c.Store.ReadProject(req.Context(), vars["name"])
	// ctx := req.Context()
	// file, _, err := req.FormFile("modelFile")
	// modelName := req.FormValue("modelName")
	return preen.ComposeModels(
		siteModel,
		ProjectModel{
			Project: selected,
		},
	)
}
