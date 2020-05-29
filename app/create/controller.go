package create

import (
	"mvc/app"
	"mvc/preen"
	"mvc/store"
	"net/http"
)

type CreateModel struct {
	ErrorMessage string
}

type CreateController struct {
	DB *store.Storage
}

func (c CreateController) Path() string {
	return "create"
}

func (c CreateController) AuthRequired() bool {
	return true
}

func (c CreateController) Get(req *http.Request) interface{} {
	return app.SiteModel{Models: c.DB.GetModels()}
}

func (c CreateController) Post(req *http.Request) interface{} {
	// file, handler, err := req.FormFile("modelFile")
	modelName := req.FormValue("modelName")

	// if err != nil {
	// 	return preen.ComposeModels(app.SiteModel{Models: c.DB.GetModels()}, CreateModel{
	// 		ErrorMessage: err.Error(),
	// 	})
	// }

	// defer file.Close()

	// buf := new(bytes.Buffer)
	// buf.ReadFrom(file)
	// content := buf.String()

	c.DB.AddModel(modelName)

	return preen.Redirect{URL: "/models/" + modelName}
}
