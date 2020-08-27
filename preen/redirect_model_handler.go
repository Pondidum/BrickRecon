package preen

import (
	"brickrecon/util"
	"context"
	"net/http"

	"github.com/honeycombio/beeline-go"
)

type controllerRedirect struct {
	controller string
	parameters map[string]interface{}
}

func ControllerRedirect(controller string, parameters ...string) interface{} {

	r := controllerRedirect{
		controller: controller,
		parameters: map[string]interface{}{},
	}

	for i := 0; i < len(parameters); i += 2 {
		key := util.Strval(parameters[i])
		value := parameters[i+1]

		r.parameters[key] = value
	}

	return r
}

type ControllerRedirectModelHandler struct {
	linker ControllerLinker
}

func NewControllerRedirectModelHandler(linker ControllerLinker) *ControllerRedirectModelHandler {

	return &ControllerRedirectModelHandler{
		linker: linker,
	}
}

func (mh ControllerRedirectModelHandler) CanHandle(ctx context.Context, model interface{}) bool {
	_, isRedirect := model.(controllerRedirect)
	beeline.AddField(ctx, "preen.is_controller_redirect", isRedirect)

	return isRedirect
}

func (mh ControllerRedirectModelHandler) Handle(ctx context.Context, ctl Controller, req *http.Request, res http.ResponseWriter, model interface{}) bool {

	redirect, isRedirect := model.(controllerRedirect)
	beeline.AddField(ctx, "preen.is_controller_redirect", isRedirect)

	if !isRedirect {
		return false
	}

	url := mh.linker(redirect.controller, redirect.parameters)

	beeline.AddField(ctx, "preen.redirect_url", url)
	http.Redirect(res, req, url, http.StatusSeeOther)

	return true
}
