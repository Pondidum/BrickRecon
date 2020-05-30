package eventstore

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
)

type Projector func(state interface{}, event interface{}) interface{}

type Projection struct {
	path            string
	initialiseState Initialiser
	project         Projector
}

func NewProjection(root string, name string, initialiseState Initialiser, project Projector) Projection {
	filepath := path.Join(root, name+".json")

	return Projection{
		path:            filepath,
		initialiseState: initialiseState,
		project:         project,
	}
}

func (p *Projection) ReadView(view interface{}) error {

	content, err := ioutil.ReadFile(p.path)

	if err != nil {
		return err
	}

	return json.Unmarshal(content, view)
}

func (p *Projection) Project(events []interface{}) error {

	state := p.initialiseState()
	err := p.ReadView(state)

	if err != nil && !os.IsNotExist(err) {
		return err
	}

	for _, e := range events {
		state = p.project(state, e)
	}

	viewBytes, err := json.Marshal(state)

	if err != nil {
		return err
	}

	return ioutil.WriteFile(p.path, viewBytes, 0666)
}
