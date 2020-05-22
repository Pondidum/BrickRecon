package command

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

type ServeCommand struct {
	Meta
}

func (c *ServeCommand) Help() string {
	return ""
}

func (c *ServeCommand) Synopsis() string {
	return "Starts the server"
}

func (c *ServeCommand) Name() string {
	return "serve"
}

func (c *ServeCommand) Run(_ []string) int {

	// pack the templates somehow
	content, _ := ioutil.ReadFile("./app/index.html")

	layout, err := template.New("layout").Parse(string(content))

	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	r := mux.NewRouter()

	templates := map[string]*template.Template{}
	err = addShared(layout, "./app/_shared")

	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	addArea(templates, "dashboard")
	addArea(templates, "create")

	r.PathPrefix("/static").Handler(http.StripPrefix("/static", http.FileServer(http.Dir("./app/static/"))))

	r.Handle("/favicon.ico", http.NotFoundHandler())
	r.HandleFunc("/{area}", func(w http.ResponseWriter, req *http.Request) {

		c.UI.Info(req.URL.String())

		clone, _ := layout.Clone()
		vars := mux.Vars(req)

		area := vars["area"]

		clone.AddParseTree("content", templates[area].Tree)

		var buffer bytes.Buffer
		err := clone.Execute(&buffer, nil)

		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte(err.Error()))
		}

		w.Write(buffer.Bytes())
	})

	c.UI.Info("Listening on 127.0.0.1:3000")
	http.ListenAndServe("127.0.0.1:3000", r)

	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	return 0
}

func addShared(templates *template.Template, pathPrefix string) error {

	files, err := ioutil.ReadDir(pathPrefix)

	if err != nil {
		return err
	}

	for _, file := range files {

		name := strings.TrimSuffix(file.Name(), path.Ext(file.Name()))
		fmt.Printf("Loading %s", name)

		content, err := ioutil.ReadFile(path.Join(pathPrefix, file.Name()))

		if err != nil {
			return err
		}

		_, err = templates.New(name).Parse(string(content))

		if err != nil {
			return err
		}

		// templates.New()
		// templates[name] = tpl
	}

	return nil
}

func addArea(templates map[string]*template.Template, name string) error {

	content, _ := ioutil.ReadFile("./app/" + name + "/index.html")
	tpl, err := template.New(name).Parse(string(content))

	if err != nil {
		return err
	}

	templates[name] = tpl

	return nil
}
