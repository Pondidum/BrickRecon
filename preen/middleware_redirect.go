package preen

import (
	"brickrecon/util"
	"context"
	"net/http"

	"github.com/honeycombio/beeline-go"
)

type redirector interface {
	RedirectUrl(ctx context.Context, c *MiddlewareContext) string
}

type controllerRedirect struct {
	controller string
	parameters map[string]interface{}
}

func (cr *controllerRedirect) RedirectUrl(ctx context.Context, c *MiddlewareContext) string {

	beeline.AddField(ctx, "preen.is_controller_redirect", true)
	return c.ControllerLink(cr.controller, cr.parameters)
}

type urlRedirect struct {
	Url string
}

func (cr *urlRedirect) RedirectUrl(ctx context.Context, c *MiddlewareContext) string {

	beeline.AddField(ctx, "preen.is_url_redirect", true)
	return cr.Url
}

func UrlRedirect(url string) interface{} {
	return &urlRedirect{
		Url: url,
	}
}

func ControllerRedirect(controller string, parameters ...string) interface{} {

	r := &controllerRedirect{
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

	if redirect, isRedirect := c.Model.(redirector); isRedirect {

		url := redirect.RedirectUrl(ctx, c)
		beeline.AddField(ctx, "preen.redirect_url", url)
		http.Redirect(response, request, url, http.StatusSeeOther)

		return false
	}

	return true
}

func isControllerRedirect(ctx context.Context, c *MiddlewareContext) (string, bool) {

	redirect, isControllerRedirect := c.Model.(controllerRedirect)
	beeline.AddField(ctx, "preen.is_controller_redirect", isControllerRedirect)

	if !isControllerRedirect {
		return "", false
	}

	url := c.ControllerLink(redirect.controller, redirect.parameters)

	return url, true
}
