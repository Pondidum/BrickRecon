package app

import (
	"brickrecon/lego"
	"brickrecon/preen"
	"net/http"
)

type KitImportController struct {
	Store *AppStore
}

func (c KitImportController) Views() []string {
	return []string{
		"kit_import_index.html",
	}
}

func (c KitImportController) Path() string {
	return "kit/import"
}

func (c KitImportController) AuthRequired() bool {
	return true
}

func (c KitImportController) Get(pc *preen.PreenContext, req *http.Request) interface{} {
	return nil
}

func (c KitImportController) Post(pc *preen.PreenContext, req *http.Request) interface{} {
	ctx := req.Context()
	kitNumber := req.FormValue("kitNumber")

	_, err := ImportKit(ctx, c.Store, lego.KitNumber(kitNumber))

	if err != nil {
		return pc.Error(err)
	}

	return preen.ControllerRedirect("kit", "kitnumber", kitNumber)
}
