package preen

import (
	"brickrecon/util"
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

func RedirectMiddleware(c *MiddlewareContext, request *http.Request, response http.ResponseWriter) bool {
	ctx := request.Context()

	redirect, isRedirect := c.Model.(controllerRedirect)
	beeline.AddField(ctx, "preen.is_controller_redirect", isRedirect)

	if !isRedirect {
		return true
	}

	url := c.ControllerLink(redirect.controller, redirect.parameters)

	beeline.AddField(ctx, "preen.redirect_url", url)
	http.Redirect(response, request, url, http.StatusSeeOther)

	return false

}
