package app

import (
	"brickrecon/lego"
	"brickrecon/lego/projections/all_kits"
	"brickrecon/preen"
	"net/http"
)

type KitModel struct {
	Kit *all_kits.KitView
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

func (c KitController) Get(pc *preen.PreenContext, req *http.Request) interface{} {

	selected, _ := c.Store.ReadKitView(req.Context(), lego.KitNumber(pc.RouteValue("kitnumber")))

	return KitModel{
		Kit: selected,
	}
}
