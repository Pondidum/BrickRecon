package preen

import "net/http"

type Controller interface {
	Path() string
}

type Getable interface {
	Get(req *http.Request) interface{}
}

type Postable interface {
	Post(req *http.Request) interface{}
}

type Redirect struct {
	URL string
}
