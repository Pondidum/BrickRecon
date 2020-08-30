package preen

import (
	"net/http"
)

type Middleware func(c *MiddlewareContext, request *http.Request, response http.ResponseWriter) bool

type MiddlewareContext struct {
	Controller    Controller
	ModelHandlers []ModelHandler
	Model         interface{}
}

func (mc *MiddlewareContext) AuthRequired() bool {
	if auth, ok := mc.Controller.(Auth); ok {
		return auth.AuthRequired()
	}

	return false
}
