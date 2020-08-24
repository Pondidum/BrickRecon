package preen

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"io/ioutil"
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
	viewRoot      string
	controllers   []Controller
	templateTypes map[string]bool
	auth          basicAuth
	getSiteModel  func(ctx context.Context) interface{}
	layout        *template.Template
	templates     map[string]*template.Template
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
		viewRoot:      pc.ApplicationRoot,
		controllers:   pc.Controllers,
		templateTypes: map[string]bool{},
		templates:     map[string]*template.Template{},
		auth:          BasicAuthMiddleware(AuthOptions{User: "test", Password: "testing"}),
		getSiteModel:  pc.GetSiteModel,
	}

	p.modelHandlers = []ModelHandler{
		&RedirectModelHandler{},
		&RenderModelHandler{
			getSiteModel: p.getSiteModel,
			render:       p.view,
		},
	}

	for _, ext := range pc.TemplateTypes {
		p.templateTypes[ext] = true
	}

	if err := p.loadViews(); err != nil {
		return p, err
	}

	if err := p.loadKnownTemplates(path.Join(pc.ApplicationRoot, "_shared")); err != nil {
		return p, err
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

func (p *Preen) loadViews() error {
	p.layout = template.New("layout").Funcs(TemplateFuncDefinitions())

	for _, c := range p.controllers {

		ctl, isController := c.(Controller)

		if !isController {
			return fmt.Errorf("%T is not a valid Controller", c)
		}

		templates, err := parseController(p.viewRoot, p.layout, ctl)
		if err != nil {
			return err
		}

		for name, tpl := range templates {
			p.templates[name] = tpl
		}
	}

	return nil
}

func parseController(viewRoot string, parentTemplate *template.Template, c Controller) (map[string]*template.Template, error) {
	templates := map[string]*template.Template{}

	for _, viewFilename := range c.Views() {

		content, err := ioutil.ReadFile(path.Join(viewRoot, viewFilename))
		if err != nil {
			return nil, err
		}

		viewPath := strings.TrimPrefix(viewFilename, controllerName(c)+"_")
		viewPath = strings.TrimSuffix(viewPath, "index.html")
		viewPath = strings.TrimSuffix(viewPath, ".html")
		viewPath = getViewName(c) + "/" + viewPath
		viewPath = strings.Trim(viewPath, "/")

		tpl := parentTemplate
		if viewPath != "" {
			tpl = parentTemplate.New(viewPath)
		}

		_, err = tpl.Parse(string(content))
		if err != nil {
			return nil, err
		}

		templates[viewPath] = tpl
	}

	return templates, nil
}

func (p *Preen) loadKnownTemplates(dir string) error {

	entries, err := ioutil.ReadDir(dir)

	if err != nil {
		return err
	}

	for _, entry := range entries {
		currentPath := path.Join(dir, entry.Name())

		if entry.IsDir() == false {

			ext := path.Ext(entry.Name())

			if !p.templateTypes[ext] {
				continue
			}

			content, err := ioutil.ReadFile(currentPath)

			if err != nil {
				return err
			}

			name := templateName("_shared", strings.TrimPrefix(currentPath, p.viewRoot+"/"))
			tpl, err := p.layout.New(name).Parse(string(content))

			if err != nil {
				return err
			}

			p.templates[name] = tpl
		} else {
			err := p.loadKnownTemplates(currentPath)

			if err != nil {
				return err
			}
		}
	}

	return nil
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

	if get, ok := c.(Getable); ok {

		r.HandleFunc("/"+ctl.Path(), func(w http.ResponseWriter, req *http.Request) {
			render(w, req, get.Get(req))
		}).Methods("GET")

	}

	if post, ok := c.(Postable); ok {

		postChain := p.auth.Wrap(render)

		r.HandleFunc("/"+ctl.Path(), func(w http.ResponseWriter, req *http.Request) {
			postChain(w, req, post.Post(req))
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

			postChain(w, req, handler(req))
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

func (p *Preen) view(w http.ResponseWriter, req *http.Request, viewName string, model interface{}) {

	clone, _ := p.layout.Clone()

	clone.Funcs(TemplateFuncs(req))

	if tpl, found := p.templates[viewName]; viewName != "" && found {
		clone.AddParseTree("content", tpl.Tree)
	} else {
		clone.New("content").Parse("")
	}

	var buffer bytes.Buffer
	err := clone.Execute(&buffer, ComposeModels(model))

	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}

	w.Write(buffer.Bytes())
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
