package preen

import (
	"context"
	"net/http"
	"reflect"
	"regexp"
	"strings"

	"github.com/gorilla/mux"
)

type ViewMiddleware func(http.ResponseWriter, *http.Request, interface{})

type Preen struct {
	viewRoot    string
	controllers []Controller
	linker      ControllerLinker
	pipeline    []Middleware
}

type PreenConfig struct {
	ApplicationRoot string

	Controllers []Controller

	TemplateTypes []string

	GetSiteModel func(ctx context.Context) interface{}
}

var defaultConfig PreenConfig = PreenConfig{
	TemplateTypes: []string{".html", ".svg"},
}

func NewPreen(pc PreenConfig) (*Preen, error) {

	if pc.TemplateTypes == nil {
		pc.TemplateTypes = defaultConfig.TemplateTypes
	}

	linker := NewControllerLinker(pc.Controllers)

	renderer, err := NewRenderMiddleware(
		pc.GetSiteModel,
		pc.ApplicationRoot,
		pc.Controllers,
		pc.TemplateTypes,
		linker)

	if err != nil {
		return nil, err
	}

	p := &Preen{
		viewRoot:    pc.ApplicationRoot,
		controllers: pc.Controllers,
		linker:      linker,
		pipeline: []Middleware{
			NewBasicAuthMiddlware("test", "testing", "Bricks").Middleware,
			RedirectMiddleware,
			TurbolinksMiddleware,
			renderer.Middleware,
		},
	}

	return p, nil
}

func (p *Preen) Apply(r *mux.Router) {

	p.HandleStaticAssets(r)

	r.Handle("/favicon.ico", http.NotFoundHandler())

	for _, ctl := range p.controllers {
		p.registerController(r, ctl)
	}
}

func (p *Preen) runMiddleware(mc *MiddlewareContext, req *http.Request, res http.ResponseWriter) {

	for _, mw := range p.pipeline {
		if mw(mc, req, res) == false {
			return
		}
	}
}

func (p *Preen) registerController(r *mux.Router, ctl Controller) error {

	controllerContext := &PreenContext{
		LinkToController: p.linker,
	}

	middlewareContext := &MiddlewareContext{
		ControllerLink: p.linker,
		Controller:     ctl,
	}

	if get, ok := ctl.(Getable); ok {

		r.HandleFunc("/"+ctl.Path(), func(w http.ResponseWriter, req *http.Request) {

			controllerContext.request = req
			middlewareContext.Model = get.Get(controllerContext, req)
			p.runMiddleware(middlewareContext, req, w)

		}).Methods("GET")

	}

	if post, ok := ctl.(Postable); ok {

		r.HandleFunc("/"+ctl.Path(), func(w http.ResponseWriter, req *http.Request) {
			controllerContext.request = req
			middlewareContext.Model = post.Post(controllerContext, req)

			p.runMiddleware(middlewareContext, req, w)

		}).Methods("POST")

	}

	if postActions, ok := ctl.(PostActions); ok {

		allActions := postActions.PostActions()

		r.HandleFunc("/"+ctl.Path(), func(w http.ResponseWriter, req *http.Request) {
			controllerContext.request = req

			action, err := getAction(controllerContext)

			if err != nil {
				middlewareContext.Model = controllerContext.Error(err)
				p.runMiddleware(middlewareContext, req, w)
				return
			}

			handler, found := allActions[action]
			if !found {
				middlewareContext.Model = controllerContext.ErrorS("No action found called " + action)
				p.runMiddleware(middlewareContext, req, w)
				return
			}

			middlewareContext.Model = handler(controllerContext, req)
			p.runMiddleware(middlewareContext, req, w)

		}).Methods("POST")

	}

	return nil
}

func (p *Preen) HandleStaticAssets(r *mux.Router) {
	r.PathPrefix("/static").Handler(http.StripPrefix("/static", http.FileServer(http.Dir("./app/static/"))))
}

func controllerName(c Controller) string {
	typeName := reflect.TypeOf(c).Elem().Name()

	name := slither(typeName)
	name = strings.ToLower(name)
	name = strings.TrimSuffix(name, "_controller")

	return name
}

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

func slither(input string) string {

	snake := matchFirstCap.ReplaceAllString(input, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

func getAction(pc *PreenContext) (string, error) {
	var pm postActions
	if err := pc.PostModel(&pm); err != nil {
		return "", err
	}

	return pm.Action, nil
}

type postActions struct {
	Action string
}
