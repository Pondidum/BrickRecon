package preen

import "net/http"

func RenderMiddlware(c *MiddlewareContext, request *http.Request, response http.ResponseWriter) bool {

	ctx := request.Context()

	for _, md := range c.ModelHandlers {

		if md.CanHandle(ctx, c.Model) == false {
			continue
		}

		md.Handle(ctx, c.Controller, request, response, c.Model)
		return true
	}

	return true
}
