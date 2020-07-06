package app

import (
	"brickrecon/lego"
	"brickrecon/preen"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	uuid "github.com/satori/go.uuid"
)

type ProjectModel struct {
	Project *ProjectWithKit
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

	projectName := lego.ProjectName(vars["name"])
	kitNumber := lego.KitNumber(req.URL.Query().Get("kit"))

	project, _ := c.Store.ReadProject(req.Context(), projectName)
	kit := project.Kits[kitNumber]

	siteModel := c.Store.SiteModel(req.Context())

	return preen.ComposeModels(
		siteModel,
		ProjectModel{
			Project: applyKit(project, kit),
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
			Project: applyKit(selected, map[lego.PartKey]int{}),
		},
	)
}

type postModel struct {
	Part     lego.LDrawPart
	Colour   lego.BrickLinkColour
	Quantity int
	Action   string
}

func applyKit(project *lego.ProjectView, kit map[lego.PartKey]int) *ProjectWithKit {

	parts := make([]PartWithKitPart, len(project.Parts))

	for i, p := range project.Parts {

		part := PartWithKitPart{
			ProjectPartView: p,
			TotalInventory:  p.Inventory,
		}

		if quantity, found := kit[lego.CreatePartKey(p.ID, p.ColourID)]; found {
			part.KitQuantity = quantity
			part.TotalInventory += quantity
		}

		parts[i] = part
	}

	return &ProjectWithKit{
		ID:    project.ID,
		Name:  project.Name,
		Parts: parts,
	}
}

type ProjectWithKit struct {
	ID   uuid.UUID
	Name lego.ProjectName

	Parts []PartWithKitPart
}

type PartWithKitPart struct {
	*lego.ProjectPartView

	KitQuantity    int
	TotalInventory int
}
