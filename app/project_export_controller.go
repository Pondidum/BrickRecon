package app

import (
	"brickrecon/lego"
	"brickrecon/preen"
	"net/http"
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

func (c ProjectExportController) Get(pc *preen.PreenContext, req *http.Request) interface{} {

	projectName := lego.ProjectName(pc.RouteValue("name"))
	project, _ := c.Store.ReadProject(req.Context(), projectName)

	return ProjectExportModel{
		WantedList:  project.BrickLinkXml,
		ProjectName: projectName,
	}

}
