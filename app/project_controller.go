package app

import (
	"brickrecon/lego"
	"brickrecon/preen"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
)

type ProjectModel struct {
	Project *lego.ProjectView
}

type ProjectController struct {
	Store *AppStore
}

func (c ProjectController) Views() []string {
	return []string{
		"project_index.html",
		"project_quantity.html",
	}
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
	selected, _ := c.Store.ReadProject(req.Context(), lego.ProjectName(vars["name"]))

	return preen.ComposeModels(
		siteModel,
		ProjectModel{
			Project: selected,
		},
	)
}

func (c ProjectController) Post(req *http.Request) interface{} {
	ctx := req.Context()
	siteModel := c.Store.SiteModel(ctx)

	if err := req.ParseForm(); err != nil {
		return preen.ComposeModels(siteModel, preen.ErrorModel(err))
	}

	decoder := schema.NewDecoder()

	var pm postModel
	if err := decoder.Decode(&pm, req.PostForm); err != nil {
		return preen.ComposeModels(siteModel, preen.ErrorModel(err))
	}

	vars := mux.Vars(req)
	projectName := lego.ProjectName(vars["name"])
	selected, _ := c.Store.ReadProject(ctx, projectName)

	project := lego.BlankProject()
	if err := c.Store.EventStore.LoadAggregate(ctx, selected.ID, project); err != nil {
		return preen.ComposeModels(siteModel, preen.ErrorModel(err))
	}

	if pm.Action == "increase" {
		if err := project.AddInventory(pm.Part, pm.Colour, pm.Quantity); err != nil {
			return preen.ComposeModels(siteModel, preen.ErrorModel(err))
		}
	}

	if pm.Action == "decrease" {
		if err := project.RemoveInventory(pm.Part, pm.Colour, pm.Quantity); err != nil {
			return preen.ComposeModels(siteModel, preen.ErrorModel(err))
		}
	}

	if err := c.Store.Save(ctx, project); err != nil {
		return preen.ComposeModels(siteModel, preen.ErrorModel(err))
	}

	selected, _ = c.Store.ReadProject(ctx, projectName)

	return preen.ComposeModels(
		siteModel,
		ProjectModel{
			Project: selected,
		},
	)
}

type postModel struct {
	Part     lego.LDrawPart
	Colour   lego.BrickLinkColour
	Quantity int
	Action   string
}
