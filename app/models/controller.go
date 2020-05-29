package models

import (
	"mvc/app"
	"mvc/store"
	"net/http"
)

type ModelsController struct {
	DB *store.Storage
}

func (c ModelsController) Path() string {
	return "models"
}

func (c ModelsController) Get(req *http.Request) interface{} {
	return app.SiteModel{AllModels: c.DB.GetModelNames()}
}
