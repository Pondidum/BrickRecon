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
	return SiteModel{AllModels: c.DB.GetModelNames()}
}
