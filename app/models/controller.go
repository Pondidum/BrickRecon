package models

import (
	"mvc/app"
	"net/http"
)

type ModelsController struct {
	Store *app.AppStore
}

func (c ModelsController) Path() string {
	return "models"
}

func (c ModelsController) Get(req *http.Request) interface{} {
	return c.Store.SiteModel()
}
