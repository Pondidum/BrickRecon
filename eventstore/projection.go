package eventstore

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
)

type Projector func(state interface{}, event Record) interface{}

type Projection struct {
	path            string
	initialiseState Initialiser
	projector       Projector
}

func NewProjection(root string, name string, initialiseState Initialiser, projector Projector) Projection {
	filepath := path.Join(root, name+".json")

	return Projection{
		path:            filepath,
		initialiseState: initialiseState,
		projector:       projector,
	}
}

func (p *Projection) ReadView(view interface{}) error {

	content, err := ioutil.ReadFile(p.path)

	if err != nil {
		return err
	}

	return json.Unmarshal(content, view)
}

func (p *Projection) Project(records []Record) error {

	state := p.initialiseState()
	err := p.ReadView(state)

	if err != nil && !os.IsNotExist(err) {
		return err
	}

	for _, r := range records {
		state = p.projector(state, r)
	}

	viewBytes, err := json.Marshal(state)

	if err != nil {
		return err
	}

	return ioutil.WriteFile(p.path, viewBytes, 0666)
}
