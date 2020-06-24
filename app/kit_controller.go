package app

import (
	"brickrecon/lego"
	"brickrecon/preen"
	"net/http"

	"github.com/gorilla/mux"
)

type KitModel struct {
	Kit *lego.KitView
}

type KitController struct {
	Store *AppStore
}

func (c KitController) Views() []string {
	return []string{
		"kit_index.html",
	}
}

func (c KitController) Path() string {
	return "kit/{kitnumber}"
}

func (c KitController) View() string {
	return "kit"
}

func (c KitController) Get(req *http.Request) interface{} {

	vars := mux.Vars(req)

	siteModel := c.Store.SiteModel(req.Context())
	selected, _ := c.Store.ReadKit(req.Context(), vars["kitnumber"])

	return preen.ComposeModels(
		siteModel,
		KitModel{
			Kit: selected,
		},
	)
}
