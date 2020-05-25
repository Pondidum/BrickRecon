package app

import (
	"net/http"
)

type AppController struct{}

func (c AppController) Path() string {
	return ""
}

func (c AppController) Get(req *http.Request) interface{} {
	return SiteModel{Models: []string{"one", "two", "three"}}
}
