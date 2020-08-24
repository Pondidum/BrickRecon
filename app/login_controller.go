package app

import (
	"brickrecon/preen"
	"net/http"
)

type LoginController struct {
}

func (c LoginController) Views() []string {
	return []string{}
}

func (c LoginController) AuthRequired() bool {
	return true
}

func (c LoginController) Path() string {
	return "login"
}

func (c LoginController) Get(req *http.Request) interface{} {
	return preen.ControllerRedirect("root")
}
