package preen

import "net/http"

type Controller interface {
	Path() string
	Views() []string
}

type Getable interface {
	Get(pc *PreenContext, req *http.Request) interface{}
}

type Postable interface {
	Post(pc *PreenContext, req *http.Request) interface{}
}

type PostActionMap map[string]func(pc *PreenContext, req *http.Request) interface{}

type PostActions interface {
	PostActions() PostActionMap // map[string]func(pc *PreenContext, req *http.Request) interface{}
}

type CustomViewName interface {
	View() string
}

type Auth interface {
	AuthRequired() bool
}
