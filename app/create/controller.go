package create

import (
	"mvc/app"
	"mvc/lego"
	"mvc/preen"
	"net/http"
)

type CreateController struct {
	Store *app.AppStore
}

func (c CreateController) Path() string {
	return "create"
}

func (c CreateController) AuthRequired() bool {
	return true
}

func (c CreateController) Get(req *http.Request) interface{} {
	return c.Store.SiteModel()
}

func (c CreateController) Post(req *http.Request) interface{} {
	file, _, err := req.FormFile("modelFile")
	modelName := req.FormValue("modelName")

	if err != nil {
		return preen.ComposeModels(c.Store.SiteModel(), preen.ErrorModel(err))
	}

	defer file.Close()

	parts, err := lego.ReadPartsList(file)

	if err != nil {
		return preen.ComposeModels(c.Store.SiteModel(), preen.ErrorModel(err))
	}

	project := lego.NewProject(modelName, parts)
	if err := c.Store.Save(project); err != nil {
		return preen.ComposeModels(c.Store.SiteModel(), preen.ErrorModel(err))
	}

	return preen.Redirect{URL: "/project/" + modelName}
}
