package app

import (
	"brickrecon/lego"
	"brickrecon/lego/projections/all_projects"
	"brickrecon/preen"
	"brickrecon/stud_io"
	"fmt"
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
		"project_kits.html",
		"project_events.html",
	}
}

func (c ProjectController) Path() string {
	return "project/{name}"
}

func (c ProjectController) View() string {
	return "project"
}

func (c ProjectController) Get(pc *preen.PreenContext, req *http.Request) interface{} {

	return ProjectModel{
		Project: projectWithKit(c.Store, pc, req),
	}

}

func (c ProjectController) PostActions() preen.PostActionMap {
	return map[string]func(pc *preen.PreenContext, req *http.Request) interface{}{
		"increase":     c.increaseQuantity,
		"decrease":     c.decreaseQuantity,
		"applykit":     c.applyKit,
		"exportWanted": c.exportWanted,
	}
}

func (c ProjectController) increaseQuantity(pc *preen.PreenContext, req *http.Request) interface{} {
	ctx := req.Context()

	project, err := c.projectAggregate(req)
	if err != nil {
		return preen.ErrorModel(err)
	}

	var pm quantityModel
	if err := preen.DecodePostForm(req.PostForm, &pm); err != nil {
		return preen.ErrorModel(err)
	}

	if err := project.AddInventory(pm.Part, pm.Colour, pm.Quantity); err != nil {
		return preen.ErrorModel(err)
	}

	if err := c.Store.Save(ctx, project); err != nil {
		return preen.ErrorModel(err)
	}

	return ProjectModel{
		Project: projectWithKit(c.Store, pc, req),
	}

}

func (c ProjectController) decreaseQuantity(pc *preen.PreenContext, req *http.Request) interface{} {
	ctx := req.Context()

	project, err := c.projectAggregate(req)
	if err != nil {
		return preen.ErrorModel(err)
	}

	var pm quantityModel
	if err := preen.DecodePostForm(req.PostForm, &pm); err != nil {
		return preen.ErrorModel(err)
	}

	if err := project.RemoveInventory(pm.Part, pm.Colour, pm.Quantity); err != nil {
		return preen.ErrorModel(err)
	}

	if err := c.Store.Save(ctx, project); err != nil {
		return preen.ErrorModel(err)
	}

	return ProjectModel{
		Project: projectWithKit(c.Store, pc, req),
	}

}

func (c ProjectController) applyKit(pc *preen.PreenContext, req *http.Request) interface{} {
	ctx := req.Context()

	vars := mux.Vars(req)

	projectName := lego.ProjectName(vars["name"])
	kitNumber := lego.KitNumber(req.URL.Query().Get("kit"))

	projectView, _ := c.Store.ReadProject(ctx, projectName)

	project := lego.BlankProject()
	if err := c.Store.EventStore.LoadAggregate(ctx, projectView.ID, project); err != nil {
		return preen.ErrorModel(err)
	}

	kit := projectView.Kits[kitNumber]

	project.AddKitContents(kit.Number, kit.Name, kitPartQuantities(kit.Parts))

	if err := c.Store.Save(ctx, project); err != nil {
		return preen.ErrorModel(err)
	}

	return ProjectModel{
		Project: projectWithKit(c.Store, pc, req),
	}

}

func (c ProjectController) exportWanted(pc *preen.PreenContext, req *http.Request) interface{} {
	ctx := req.Context()

	project, err := c.projectAggregate(req)
	if err != nil {
		return preen.ErrorModel(err)
	}

	exporter := &stud_io.WantedListXmlExporter{}

	if _, err := project.ExportWantedList(exporter); err != nil {
		return preen.ErrorModel(err)
	}

	if err := c.Store.Save(ctx, project); err != nil {
		return preen.ErrorModel(err)
	}

	return preen.ControllerRedirect("project_export", "name", string(project.Name))
}

func kitPartQuantities(quantities map[all_projects.PartKey]int) []lego.PartQuantity {

	parts := make([]lego.PartQuantity, len(quantities))
	i := 0
	for key, q := range quantities {
		part, colour := all_projects.ParseKey(key)

		parts[i] = lego.PartQuantity{PartID: part, ColourID: colour, Quantity: q}
		i++
	}

	return parts
}

func (c ProjectController) projectAggregate(req *http.Request) (*lego.Project, error) {

	ctx := req.Context()
	vars := mux.Vars(req)

	projectName := lego.ProjectName(vars["name"])
	selected, _ := c.Store.ReadProject(ctx, projectName)

	project := lego.BlankProject()
	if err := c.Store.EventStore.LoadAggregate(ctx, selected.ID, project); err != nil {
		return nil, err
	}

	return project, nil
}

type quantityModel struct {
	Part     lego.LDrawPart
	Colour   lego.BrickLinkColour
	Quantity int
}

func projectWithKit(store *AppStore, pc *preen.PreenContext, req *http.Request) *ProjectWithKit {
	vars := mux.Vars(req)

	projectName := lego.ProjectName(vars["name"])
	kitNumber := lego.KitNumber(req.URL.Query().Get("kit"))

	project, _ := store.ReadProject(req.Context(), projectName)
	kit := project.Kits[kitNumber]

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
		ID:     project.ID,
		Name:   project.Name,
		Parts:  parts,
		Kits:   project.Kits,
		Events: withLinks(pc, project.Events),
	}
}

func withLinks(pc *preen.PreenContext, events []*all_projects.EventDescription) []*all_projects.EventDescription {

	for _, event := range events {
		if event.Type == "KitAddedToProject" {
			event.Description = fmt.Sprintf(
				`Kit <a href="%s">%s</a> added`,
				pc.LinkToController("kit", event.Additional),
				event.Additional["KitName"],
			)
		}
	}

	return events
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

	Parts  []PartWithKitPart
	Kits   map[lego.KitNumber]all_projects.KitView
	Events []*all_projects.EventDescription
}

type PartWithKitPart struct {
	*all_projects.ProjectPartView

	KitQuantity    int
	TotalInventory int
}
