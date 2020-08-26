package app

import (
	"brickrecon/lego"
	"net/http"

	"github.com/gorilla/mux"
)

type ProjectExportModel struct {
	WantedList  string
	ProjectName lego.ProjectName
}

type ProjectExportController struct {
	Store *AppStore
}

func (c ProjectExportController) Views() []string {
	return []string{
		"project_export_index.html",
	}
}

func (c ProjectExportController) Path() string {
	return "project/{name}/export"
}

func (c ProjectExportController) View() string {
	return "project_export"
}

func (c ProjectExportController) Get(req *http.Request) interface{} {
	vars := mux.Vars(req)

	projectName := lego.ProjectName(vars["name"])
	project, _ := c.Store.ReadProject(req.Context(), projectName)

	return ProjectExportModel{
		WantedList:  project.BrickLinkXml,
		ProjectName: projectName,
	}

}
