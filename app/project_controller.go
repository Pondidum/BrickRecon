package app

import (
	"brickrecon/lego"
	"brickrecon/lego/projections/all_projects"
	"brickrecon/preen"
	"net/http"
	"sort"

	"github.com/gorilla/mux"
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

	return ProjectModel{
		Project: projectWithKit(c.Store, req),
	}

}

func projectWithKit(store *AppStore, req *http.Request) *ProjectWithKit {
	vars := mux.Vars(req)

	projectName := lego.ProjectName(vars["name"])
	kitNumber := lego.KitNumber(req.URL.Query().Get("kit"))

	project, _ := store.ReadProject(req.Context(), projectName)
	kit := project.Kits[kitNumber]

	return applyKit(project, kit)
}

var actions = map[string]func(project *lego.Project, req *http.Request) error{
	"increase": handleIncrease,
	"decrease": handleDecrease,
}

func handleIncrease(project *lego.Project, req *http.Request) error {

	var pm quantityModel
	if err := preen.DecodePostForm(req.PostForm, &pm); err != nil {
		return err
	}

	if err := project.AddInventory(pm.Part, pm.Colour, pm.Quantity); err != nil {
		return err
	}

	return nil
}

func handleDecrease(project *lego.Project, req *http.Request) error {

	var pm quantityModel
	if err := preen.DecodePostForm(req.PostForm, &pm); err != nil {

		return err
	}

	if err := project.RemoveInventory(pm.Part, pm.Colour, pm.Quantity); err != nil {
		return err
	}

	return nil
}

func (c ProjectController) Post(req *http.Request) interface{} {
	ctx := req.Context()

	if err := req.ParseForm(); err != nil {
		return preen.ErrorModel(err)
	}

	vars := mux.Vars(req)
	projectName := lego.ProjectName(vars["name"])
	selected, _ := c.Store.ReadProject(ctx, projectName)

	project := lego.BlankProject()
	if err := c.Store.EventStore.LoadAggregate(ctx, selected.ID, project); err != nil {
		return preen.ErrorModel(err)
	}

	action, err := getAction(req)
	if err != nil {
		return preen.ErrorModel(err)
	}

	handler, found := actions[action]
	if !found {
		return preen.ErrorModelS("No handler found for action " + action)
	}

	if err := handler(project, req); err != nil {
		return preen.ErrorModel(err)
	}

	if err := c.Store.Save(ctx, project); err != nil {
		return preen.ErrorModel(err)
	}

	return ProjectModel{
		Project: projectWithKit(c.Store, req),
	}

}

func getAction(req *http.Request) (string, error) {

	var pm postActions
	if err := preen.DecodePostForm(req.PostForm, &pm); err != nil {
		return "", err
	}

	return pm.Action, nil
}

type postActions struct {
	Action string
}

type quantityModel struct {
	Part     lego.LDrawPart
	Colour   lego.BrickLinkColour
	Quantity int
}

func applyKit(project *all_projects.ProjectView, kit all_projects.KitView) *ProjectWithKit {

	parts := make([]PartWithKitPart, len(project.Parts))

	for i, p := range project.Parts {

		part := PartWithKitPart{
			ProjectPartView: p,
			TotalInventory:  p.Inventory,
		}

		pk := all_projects.CreatePartKey(p.ID, p.ColourID)

		if quantity, found := kit.Parts[pk]; found {
			part.KitQuantity = quantity
			part.TotalInventory += quantity
		}

		parts[i] = part
	}

	if kit.Name != lego.KitName("") {
		sortByKitAddition(parts)
	} else {
		sortByPartAndColour(parts)
	}

	return &ProjectWithKit{
		ID:    project.ID,
		Name:  project.Name,
		Parts: parts,
		Kits:  project.Kits,
	}
}

func sortByKitAddition(parts []PartWithKitPart) {
	sort.Slice(parts, func(x int, y int) bool {
		l := parts[x]
		r := parts[y]

		if (l.KitQuantity > 0) == (r.KitQuantity > 0) {

			if l.ID == r.ID {
				return l.ColourID < r.ColourID
			}

			return l.ID < r.ID
		}

		return l.KitQuantity > 0
	})
}

func sortByPartAndColour(parts []PartWithKitPart) {
	sort.Slice(parts, func(x int, y int) bool {
		l := parts[x]
		r := parts[y]

		if l.ID == r.ID {
			return l.ColourID < r.ColourID
		}

		return l.ID < r.ID
	})
}

type ProjectWithKit struct {
	ID   uuid.UUID
	Name lego.ProjectName

	Parts []PartWithKitPart
	Kits  map[lego.KitNumber]all_projects.KitView
}

type PartWithKitPart struct {
	*all_projects.ProjectPartView

	KitQuantity    int
	TotalInventory int
}
