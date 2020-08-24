package preen

import (
	"context"
	"net/http"

	"github.com/honeycombio/beeline-go"
)

type RenderModelHandler struct {
	getSiteModel func(ctx context.Context) interface{}
	render       func(w http.ResponseWriter, req *http.Request, viewName string, model interface{})
}

func (mh *RenderModelHandler) CanHandle(ctx context.Context, model interface{}) bool {
	return true
}

func (mh *RenderModelHandler) Handle(ctx context.Context, ctl Controller, req *http.Request, res http.ResponseWriter, model interface{}) bool {

	siteModel := mh.getSiteModel(ctx)
	viewModel := ComposeModels(siteModel, model)
	viewName := getViewName(ctl)

	beeline.AddField(ctx, "preen.view_name", viewName)

	mh.render(res, req, viewName, viewModel)

	return true
}
