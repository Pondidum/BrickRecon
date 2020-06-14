package app

import (
	"net/http"
)

type AppController struct {
	Store *AppStore
}

func (c AppController) Path() string {
	return ""
}

func (c AppController) Get(req *http.Request) interface{} {
	return c.Store.SiteModel(req.Context())
}
