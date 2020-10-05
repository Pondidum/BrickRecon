package preen

import (
	"net/http"

	"github.com/gorilla/schema"
)

type PreenContext struct {
	request          *http.Request
	LinkToController ControllerLinker
}

func (pc *PreenContext) Redirect(url string) interface{} {
	return UrlRedirect(url)
}

func (pc *PreenContext) Error(err error) Error {
	return Error{ErrorMessage: err.Error()}
}

func (pc *PreenContext) ErrorS(err string) Error {
	return Error{ErrorMessage: err}
}

var decoder = schema.NewDecoder()

func (pc *PreenContext) PostModel(model interface{}) error {

	decoder.IgnoreUnknownKeys(true)

	if err := pc.request.ParseForm(); err != nil {
		return err
	}

	return decoder.Decode(model, pc.request.PostForm)
}
