package command

import (
	"bytes"
	"html/template"
	"io/ioutil"
	"net/http"

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
	content, _ := ioutil.ReadFile("./app/ui/layout.html")

	layout, err := template.New("layout").Parse(string(content))

	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	r := mux.NewRouter()

	templates := map[string]*template.Template{}
	addArea(templates, "dashboard")
	addArea(templates, "create")

	r.HandleFunc("/{area}", func(w http.ResponseWriter, req *http.Request) {
		clone, _ := layout.Clone()
		vars := mux.Vars(req)

		area := vars["area"]

		clone.AddParseTree("content", templates[area].Tree)

		var buffer bytes.Buffer
		clone.Execute(&buffer, nil)

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

func addArea(templates map[string]*template.Template, name string) error {

	content, _ := ioutil.ReadFile("./app/" + name + "/index.html")
	tpl, err := template.New(name).Parse(string(content))

	if err != nil {
		return err
	}

	templates[name] = tpl

	return nil
}
