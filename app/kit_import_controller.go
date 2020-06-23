package app

import (
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

func (c KitImportController) Get(req *http.Request) interface{} {
	return c.Store.SiteModel(req.Context())
}

func (c KitImportController) Post(req *http.Request) interface{} {
	// ctx := req.Context()
	kitNumber := req.FormValue("kitNumber")

	return preen.Redirect{URL: "/kit/" + kitNumber}
}
