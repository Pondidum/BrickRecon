package preen

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"path"
	"strings"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/mitchellh/mapstructure"
)

type ViewMiddleware func(http.ResponseWriter, *http.Request, interface{})

type Preen struct {
	viewRoot    string
	controllers []Controller
	auth        basicAuth

	layout    *template.Template
	templates map[string]*template.Template
}

type PreenConfig struct {
	ApplicationRoot string

	Controllers []Controller
}

func NewPreen(pc PreenConfig) (Preen, error) {
	p := Preen{
		viewRoot:    pc.ApplicationRoot,
		controllers: pc.Controllers,
		templates:   map[string]*template.Template{},
		auth:        BasicAuthMiddleware(AuthOptions{User: "test", Password: "testing"}),
	}

	if err := p.loadLayoutRoot(); err != nil {
		return p, err
	}

	if err := p.loadTemplates(pc.ApplicationRoot); err != nil {
		return p, err
	}

	return p, nil
}

func (p *Preen) Apply(r *mux.Router) {

	p.HandleStaticAssets(r)

	r.Handle("/favicon.ico", http.NotFoundHandler())
	r.Use(p.auth.UserContext)

	for _, ctl := range p.controllers {
		p.RegisterController(r, ctl)
	}
}

func (p *Preen) loadLayoutRoot() error {
	content, err := ioutil.ReadFile(path.Join(p.viewRoot, "index.html"))

	if err != nil {
		return err
	}

	layout, err := template.New("layout").Parse(string(content))

	if err != nil {
		return err
	}

	p.layout = layout

	return nil
}

func (p *Preen) loadTemplates(dir string) error {

	entries, err := ioutil.ReadDir(dir)

	if err != nil {
		return err
	}

	for _, entry := range entries {
		currentPath := path.Join(dir, entry.Name())

		if entry.IsDir() == false {

			content, err := ioutil.ReadFile(currentPath)

			if err != nil {
				return err
			}

			name := templateName(strings.TrimPrefix(currentPath, p.viewRoot+"/"))
			tpl, err := p.layout.New(name).Parse(string(content))

			if err != nil {
				return err
			}

			p.templates[name] = tpl
		} else {
			err := p.loadTemplates(currentPath)

			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (p *Preen) RegisterController(r *mux.Router, c interface{}) error {

	ctl, isController := c.(Controller)

	if !isController {
		return fmt.Errorf("%T is not a valid Controller", c)
	}

	render := func(w http.ResponseWriter, req *http.Request, model interface{}) {
		if redirect, ok := model.(Redirect); ok {
			http.Redirect(w, req, redirect.URL, http.StatusSeeOther)
		} else {
			p.view(w, req, getViewName(ctl), model)
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

		r.HandleFunc("/"+ctl.Path(), func(w http.ResponseWriter, req *http.Request) {
			render(w, req, post.Post(req))
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

	context, err := composeModel(req, model)

	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}

	clone, _ := p.layout.Clone()

	if tpl, found := p.templates[viewName]; viewName != "" && found {
		clone.AddParseTree("content", tpl.Tree)
	} else {
		clone.New("content").Parse("")
	}

	var buffer bytes.Buffer
	err = clone.Execute(&buffer, context)

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

func composeModel(req *http.Request, model interface{}) (map[string]interface{}, error) {

	user := context.Get(req, "UserInfo")

	if user == nil {
		user = UserInfo{}
	}

	site := SiteInfo{
		URL: req.Host,
	}

	ctx := map[string]interface{}{
		"_PagePath": req.URL.Path,
		"_User":     user,
		"_Site":     site,
	}

	if err := mapstructure.Decode(model, &ctx); err != nil {
		return nil, err
	}

	return ctx, nil
}

func templateName(filepath string) string {
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

type UserInfo struct {
	Name          string
	Authenticated bool
}

type SiteInfo struct {
	URL string
}
