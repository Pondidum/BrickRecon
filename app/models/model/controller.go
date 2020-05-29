package model

import (
	"mvc/app"
	"mvc/store"
	"net/http"

	"github.com/gorilla/mux"
)

type ModelController struct {
	DB *store.Storage
}

func (c ModelController) Path() string {
	return "models/{name}"
}

func (c ModelController) View() string {
	return "models/model"
}

func (c ModelController) Get(req *http.Request) interface{} {

	vars := mux.Vars(req)
	names := c.DB.GetModelNames()

	return app.SiteModel{AllModels: names, SelectedModel: c.DB.GetModel(vars["name"])}
}
