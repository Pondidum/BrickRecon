package preen

import (
	"context"
	"net/http"
)

type ModelHandler interface {
	CanHandle(ctx context.Context, model interface{}) bool
	Handle(ctx context.Context, ctl Controller, req *http.Request, res http.ResponseWriter, model interface{}) bool
}
