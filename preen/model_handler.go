package preen

import (
	"context"
	"net/http"

	"github.com/honeycombio/beeline-go"
)

type ModelHandler interface {
	CanHandle(ctx context.Context, model interface{}) bool
	Handle(ctx context.Context, ctl Controller, req *http.Request, res http.ResponseWriter, model interface{}) bool
}

// ------------------------------------------------------------------------- //

type RedirectModelHandler struct{}

func (mh *RedirectModelHandler) CanHandle(ctx context.Context, model interface{}) bool {
	_, isRedirect := model.(Redirect)
	beeline.AddField(ctx, "preen.is_redirect", isRedirect)

	return isRedirect
}

func (mh *RedirectModelHandler) Handle(ctx context.Context, ctl Controller, req *http.Request, res http.ResponseWriter, model interface{}) bool {
	redirect, isRedirect := model.(Redirect)
	beeline.AddField(ctx, "preen.is_redirect", isRedirect)

	if !isRedirect {
		return false
	}

	beeline.AddField(ctx, "preen.redirect_url", redirect.URL)
	http.Redirect(res, req, redirect.URL, http.StatusSeeOther)
	return true
}

// ------------------------------------------------------------------------- //

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
