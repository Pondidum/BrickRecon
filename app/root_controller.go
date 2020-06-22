package app

import (
	"net/http"
)

type RootController struct {
	Store *AppStore
}

func (c RootController) Views() []string {
	return []string{
		"root_index.html",
		"menu/index.html",
		"menu/project-list.html",
		"menu/sets-list.html",
	}
}

func (c RootController) Path() string {
	return ""
}

func (c RootController) Get(req *http.Request) interface{} {
	return c.Store.SiteModel(req.Context())
}
