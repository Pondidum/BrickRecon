package app

import (
	"brickrecon/bricklink"
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

	ctx := req.Context()

	projectName := lego.ProjectName(pc.RouteValue("name"))

	project, err := c.Store.ReadProject(ctx, projectName)
	if err != nil {
		return pc.Error(err)
	}

	// move this to a projection later
	wantedParts := []*bricklink.WantedListPart{}
	for key, amounts := range project.Parts() {
		part, err := c.Store.ReadPart(ctx, key)
		if err != nil {
			return pc.Error(err)
		}

		wantedParts = append(wantedParts, &bricklink.WantedListPart{
			ID:        part.BrickLink.PartNumber,
			Colour:    part.BrickLink.Colour,
			Quantity:  amounts.Quantity,
			Inventory: amounts.Inventory,
		})
	}

	exporter := &bricklink.WantedListXmlExporter{}
	markup, err := exporter.Export(wantedParts)
	if err != nil {
		return pc.Error(err)
	}

	return ProjectExportModel{
		WantedList:  markup,
		ProjectName: projectName,
	}

}
