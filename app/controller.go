package app

import (
	"mvc/store"
	"net/http"
)

type AppController struct {
	DB *store.Storage
}

func (c AppController) Path() string {
	return ""
}

func (c AppController) Get(req *http.Request) interface{} {
	return SiteModel{Models: c.DB.GetModels()}
}
