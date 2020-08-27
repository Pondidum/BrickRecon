package preen

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"reflect"
	"regexp"
	"strings"

	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"github.com/mitchellh/mapstructure"
)

type ViewMiddleware func(http.ResponseWriter, *http.Request, interface{})

type Preen struct {
	viewRoot    string
	controllers []Controller

	auth          basicAuth
	modelHandlers []ModelHandler
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

func NewPreen(pc PreenConfig) (Preen, error) {

	if pc.TemplateTypes == nil {
		pc.TemplateTypes = defaultConfig.TemplateTypes
	}

	p := Preen{
		viewRoot:    pc.ApplicationRoot,
		controllers: pc.Controllers,
		auth:        BasicAuthMiddleware(AuthOptions{User: "test", Password: "testing"}),
	}

	renderer, err := NewRenderModelHandler(pc.GetSiteModel, pc.ApplicationRoot, pc.Controllers, pc.TemplateTypes)
	if err != nil {
		return p, err
	}

	p.modelHandlers = []ModelHandler{
		NewControllerRedirectModelHandler(p.controllers),
		renderer,
	}

	return p, nil
}

func (p *Preen) Apply(r *mux.Router) {

	p.HandleStaticAssets(r)

	r.Handle("/favicon.ico", http.NotFoundHandler())
	r.Use(p.auth.UserContext)

	for _, ctl := range p.controllers {
		p.registerController(r, ctl)
	}
}

func (p *Preen) registerController(r *mux.Router, c interface{}) error {

	ctl, isController := c.(Controller)

	if !isController {
		return fmt.Errorf("%T is not a valid Controller", c)
	}

	render := func(w http.ResponseWriter, req *http.Request, model interface{}) {
		ctx := req.Context()

		for _, md := range p.modelHandlers {
			if md.CanHandle(ctx, model) && md.Handle(ctx, ctl, req, w, model) {
				continue
			}
		}

	}

	if auth, ok := c.(Auth); ok {
		if auth.AuthRequired() {
			render = p.auth.Wrap(render)
		}
	}

	context := &PreenContext{}

	if get, ok := c.(Getable); ok {

		r.HandleFunc("/"+ctl.Path(), func(w http.ResponseWriter, req *http.Request) {
			render(w, req, get.Get(context, req))
		}).Methods("GET")

	}

	if post, ok := c.(Postable); ok {

		postChain := p.auth.Wrap(render)

		r.HandleFunc("/"+ctl.Path(), func(w http.ResponseWriter, req *http.Request) {
			postChain(w, req, post.Post(context, req))
		}).Methods("POST")

	}

	if postActions, ok := c.(PostActions); ok {

		postChain := p.auth.Wrap(render)
		allActions := postActions.PostActions()

		r.HandleFunc("/"+ctl.Path(), func(w http.ResponseWriter, req *http.Request) {
			action, err := getAction(req)

			if err != nil {
				render(w, req, ErrorModel(err))
				return
			}

			handler, found := allActions[action]
			if !found {
				render(w, req, ErrorModelS("No action found called "+action))
				return
			}

			postChain(w, req, handler(context, req))
		}).Methods("POST")

	}

	return nil
}

func getViewName(ctl Controller) string {

	if custom, ok := ctl.(CustomViewName); ok {
		return custom.View()
	}

	return ctl.Path()
}

func (p *Preen) HandleStaticAssets(r *mux.Router) {
	r.PathPrefix("/static").Handler(http.StripPrefix("/static", http.FileServer(http.Dir("./app/static/"))))
}

func ComposeModels(models ...interface{}) interface{} {

	result := map[string]interface{}{}

	for _, m := range models {
		mapstructure.Decode(m, &result)
	}

	return result
}

func templateName(controller string, filepath string) string {

	if strings.HasPrefix(filepath, controller+"_") {
		filepath = strings.Replace(filepath, controller+"_", controller+"/", 1)
	}

	ext := path.Ext(filepath)
	base := path.Base(filepath)

	if base == "index.html" {
		filepath = strings.TrimSuffix(filepath, base)
	}

	filepath = strings.TrimSuffix(filepath, ext)
	filepath = strings.TrimSuffix(filepath, "/")

	filepath = strings.TrimPrefix(filepath, "_shared/")

	return filepath
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

type PreenContext struct{}
