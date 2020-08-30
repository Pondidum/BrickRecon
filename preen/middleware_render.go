package preen

import (
	"bytes"
	"context"
	"html/template"
	"io/ioutil"
	"net/http"
	"path"
	"strings"

	"github.com/honeycombio/beeline-go"
	"github.com/mitchellh/mapstructure"
)

type RenderModelHandler struct {
	getSiteModel func(ctx context.Context) interface{}

	layout        *template.Template
	templates     map[string]*template.Template
	templateTypes map[string]bool
	Context       *TemplateContext
}

func NewRenderMiddleware(getSiteModel func(ctx context.Context) interface{}, viewRoot string, controllers []Controller, templateTypes []string, linker ControllerLinker) (*RenderModelHandler, error) {

	context := NewTemplateContext(linker)

	mh := &RenderModelHandler{
		getSiteModel:  getSiteModel,
		layout:        template.New("layout").Funcs(context.Functions),
		templateTypes: map[string]bool{},
		Context:       context,
	}

	if err := mh.loadViews(viewRoot, controllers); err != nil {
		return nil, err
	}

	for _, ext := range templateTypes {
		mh.templateTypes[ext] = true
	}

	if err := mh.loadKnownTemplates(viewRoot, path.Join(viewRoot, "_shared")); err != nil {
		return mh, err
	}

	return mh, nil
}

func (mh *RenderModelHandler) Middleware(c *MiddlewareContext, request *http.Request, response http.ResponseWriter) bool {

	ctx := request.Context()
	siteModel := mh.getSiteModel(ctx)
	viewModel := composeModels(siteModel, c.Model)
	viewName := getViewName(c.Controller)

	beeline.AddField(ctx, "preen.view_name", viewName)

	mh.render(c.Controller, request, response, viewName, viewModel)

	return true
}

func (mh *RenderModelHandler) loadViews(viewRoot string, controllers []Controller) error {

	templates := map[string]*template.Template{}

	for _, c := range controllers {

		controllerTemplates, err := mh.parseController(viewRoot, c)
		if err != nil {
			return err
		}

		for name, tpl := range controllerTemplates {
			templates[name] = tpl
		}
	}

	mh.templates = templates

	return nil
}

func (mh *RenderModelHandler) parseController(viewRoot string, c Controller) (map[string]*template.Template, error) {
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

		tpl := mh.layout
		if viewPath != "" {
			tpl = mh.layout.New(viewPath)
		}

		_, err = tpl.Parse(string(content))
		if err != nil {
			return nil, err
		}

		templates[viewPath] = tpl
	}

	return templates, nil
}

func (mh *RenderModelHandler) loadKnownTemplates(viewRoot string, dir string) error {

	entries, err := ioutil.ReadDir(dir)

	if err != nil {
		return err
	}

	for _, entry := range entries {
		currentPath := path.Join(dir, entry.Name())

		if entry.IsDir() == false {

			ext := path.Ext(entry.Name())

			if !mh.templateTypes[ext] {
				continue
			}

			content, err := ioutil.ReadFile(currentPath)

			if err != nil {
				return err
			}

			name := templateName("_shared", strings.TrimPrefix(currentPath, viewRoot+"/"))
			tpl, err := mh.layout.New(name).Parse(string(content))

			if err != nil {
				return err
			}

			mh.templates[name] = tpl
		} else {
			err := mh.loadKnownTemplates(viewRoot, currentPath)

			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (mh *RenderModelHandler) render(ctl Controller, req *http.Request, w http.ResponseWriter, viewName string, model interface{}) {

	clone, _ := mh.layout.Clone()

	mh.Context.Request = req

	clone.Funcs(mh.Context.Functions)

	if tpl, found := mh.templates[viewName]; viewName != "" && found {
		clone.AddParseTree("content", tpl.Tree)
	} else {
		clone.New("content").Parse("")
	}

	var buffer bytes.Buffer
	err := clone.Execute(&buffer, composeModels(model))

	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}

	w.Write(buffer.Bytes())
}

func getViewName(ctl Controller) string {

	if custom, ok := ctl.(CustomViewName); ok {
		return custom.View()
	}

	return ctl.Path()
}

func composeModels(models ...interface{}) interface{} {

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
