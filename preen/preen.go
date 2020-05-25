package preen

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"path"
	"strings"

	"github.com/gorilla/mux"
)

type Preen struct {
	viewRoot string

	layout    *template.Template
	templates map[string]*template.Template
}

func NewPreen(viewRoot string) (Preen, error) {
	p := Preen{
		viewRoot:  viewRoot,
		templates: map[string]*template.Template{},
	}

	if err := p.loadLayoutRoot(); err != nil {
		return p, err
	}

	if err := p.loadTemplates(viewRoot); err != nil {
		return p, err
	}

	return p, nil
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

	if get, ok := c.(Getable); ok {

		r.HandleFunc("/"+ctl.Path(), func(w http.ResponseWriter, req *http.Request) {
			p.View(w, ctl.Path(), get.Get(req))
		}).Methods("GET")

	}

	if post, ok := c.(Postable); ok {

		r.HandleFunc("/"+ctl.Path(), func(w http.ResponseWriter, req *http.Request) {
			p.View(w, ctl.Path(), post.Post(req))
		}).Methods("POST")

	}

	return nil
}

func (p *Preen) HandleStaticAssets(r *mux.Router) {
	r.PathPrefix("/static").Handler(http.StripPrefix("/static", http.FileServer(http.Dir("./app/static/"))))
}

func (p *Preen) View(w http.ResponseWriter, viewName string, model interface{}) {

	clone, _ := p.layout.Clone()

	if tpl, found := p.templates[viewName]; found {
		clone.AddParseTree("content", tpl.Tree)
	} else {
		clone.New("content").Parse("")
	}

	var buffer bytes.Buffer
	err := clone.Execute(&buffer, model)

	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
	}

	w.Write(buffer.Bytes())
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
