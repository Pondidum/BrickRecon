package preen

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"regexp"
	"strings"

	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
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

func (p *Preen) registerController(r *mux.Router, c interface{}) error {

	ctl, isController := c.(Controller)

	if !isController {
		return fmt.Errorf("%T is not a valid Controller", c)
	}

	context := &PreenContext{
		LinkToController: p.linker,
	}

	mc := &MiddlewareContext{
		ControllerLink: p.linker,
		Controller:     ctl,
	}

	if get, ok := c.(Getable); ok {

		r.HandleFunc("/"+ctl.Path(), func(w http.ResponseWriter, req *http.Request) {

			mc.Model = get.Get(context, req)
			p.runMiddleware(mc, req, w)

		}).Methods("GET")

	}

	if post, ok := c.(Postable); ok {

		r.HandleFunc("/"+ctl.Path(), func(w http.ResponseWriter, req *http.Request) {
			mc.Model = post.Post(context, req)

			p.runMiddleware(mc, req, w)

		}).Methods("POST")

	}

	if postActions, ok := c.(PostActions); ok {

		allActions := postActions.PostActions()

		r.HandleFunc("/"+ctl.Path(), func(w http.ResponseWriter, req *http.Request) {
			action, err := getAction(req)

			if err != nil {
				mc.Model = context.Error(err)
				p.runMiddleware(mc, req, w)
				return
			}

			handler, found := allActions[action]
			if !found {
				mc.Model = context.ErrorS("No action found called " + action)
				p.runMiddleware(mc, req, w)
				return
			}

			mc.Model = handler(context, req)
			p.runMiddleware(mc, req, w)

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

var decoder = schema.NewDecoder()

func DecodePostForm(form url.Values, model interface{}) error {

	decoder.IgnoreUnknownKeys(true)

	err := decoder.Decode(model, form)
	return err
}

func getAction(req *http.Request) (string, error) {

	if err := req.ParseForm(); err != nil {
		return "", err
	}

	var pm postActions
	if err := DecodePostForm(req.PostForm, &pm); err != nil {
		return "", err
	}

	return pm.Action, nil
}

type postActions struct {
	Action string
}
