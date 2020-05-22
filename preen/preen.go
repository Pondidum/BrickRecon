package preen

import (
	"bytes"
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

	if err := p.loadShared(); err != nil {
		return p, err
	}

	if err := p.loadAreas(); err != nil {
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

func (p *Preen) loadShared() error {
	shared := path.Join(p.viewRoot, "_shared")
	files, err := ioutil.ReadDir(shared)

	if err != nil {
		return err
	}

	for _, file := range files {

		name := strings.TrimSuffix(file.Name(), path.Ext(file.Name()))

		content, err := ioutil.ReadFile(path.Join(shared, file.Name()))

		if err != nil {
			return err
		}

		_, err = p.layout.New(name).Parse(string(content))

		if err != nil {
			return err
		}
	}

	return nil
}

func (p *Preen) loadAreas() error {

	dirs, err := ioutil.ReadDir(p.viewRoot)

	if err != nil {
		return err
	}

	for _, dir := range dirs {

		if !dir.IsDir() {
			continue
		}

		files, err := ioutil.ReadDir(path.Join(p.viewRoot, dir.Name()))

		if err != nil {
			return err
		}

		for _, file := range files {

			if path.Ext(file.Name()) == ".html" {

				templateName := strings.TrimSuffix(file.Name(), path.Ext(file.Name()))

				if templateName == "index" {
					templateName = dir.Name()
				} else {
					templateName = dir.Name() + "/" + templateName
				}

				content, _ := ioutil.ReadFile(path.Join(p.viewRoot, dir.Name(), file.Name()))

				tpl, err := p.layout.New(templateName).Parse(string(content))

				if err != nil {
					return err
				}

				p.templates[templateName] = tpl
			}
		}
	}

	return nil

}

func (p *Preen) HandleStaticAssets(r *mux.Router) {
	r.PathPrefix("/static").Handler(http.StripPrefix("/static", http.FileServer(http.Dir("./app/static/"))))
}

func (p *Preen) View(w http.ResponseWriter, req *http.Request, model interface{}) {

	clone, _ := p.layout.Clone()

	vars := mux.Vars(req)
	if area, found := vars["area"]; found {
		clone.AddParseTree("content", p.templates[area].Tree)
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
