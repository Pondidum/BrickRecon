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
		"menu/kit-list.html",
	}
}

func (c RootController) Path() string {
	return ""
}

func (c RootController) Get(req *http.Request) interface{} {
	return nil
}
