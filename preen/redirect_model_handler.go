package preen

import (
	"context"
	"net/http"
	"regexp"
	"strings"

	"github.com/honeycombio/beeline-go"
)

type controllerRedirect struct {
	controller string
	parameters map[string]string
}

func ControllerRedirect(controller string, parameters ...string) interface{} {

	r := controllerRedirect{
		controller: controller,
		parameters: map[string]string{},
	}

	for i := 0; i < len(parameters); i += 2 {
		key := strval(parameters[i])
		value := strval(parameters[i+1])

		r.parameters[key] = value
	}

	return r
}

var rx = regexp.MustCompile("{(.*?)}")

type ControllerRedirectModelHandler struct {
	controllers map[string]Controller
}

func NewControllerRedirectModelHandler(controllers []Controller) *ControllerRedirectModelHandler {

	lookup := map[string]Controller{}

	for _, ctl := range controllers {
		lookup[controllerName(ctl)] = ctl
	}

	return &ControllerRedirectModelHandler{
		controllers: lookup,
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

	toController := mh.controllers[redirect.controller]

	url := "/" + rx.ReplaceAllStringFunc(toController.Path(), func(match string) string {
		return redirect.parameters[strings.Trim(match, "{}")]
	})

	beeline.AddField(ctx, "preen.redirect_url", url)
	http.Redirect(res, req, url, http.StatusSeeOther)

	return true
}
