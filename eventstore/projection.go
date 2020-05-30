package eventstore

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
	"strconv"
)

type Projector func(state interface{}, event interface{}) interface{}

type Projection struct {
	root         string
	name         string
	initialState Initialiser
	project      Projector
}

func NewProjection(root string, name string, initialisState Initialiser, project Projector) Projection {
	return Projection{root, name, initialisState, project}
}

func (p *Projection) CheckIndex() (int, error) {

	filename := path.Join(p.root, p.name+".idx")
	contents, err := ioutil.ReadFile(filename)

	if os.IsNotExist(err) {
		return 0, nil
	}

	if err != nil {
		return 0, err
	}

	return strconv.Atoi(string(contents))
}

func (p *Projection) WriteCheckIndex(index int) error {
	filename := path.Join(p.root, p.name+".idx")
	contents := strconv.Itoa(index)

	return ioutil.WriteFile(filename, []byte(contents), 0666)
}

func (p *Projection) ReadView(view interface{}) error {

	filename := path.Join(p.root, p.name+".json")
	content, err := ioutil.ReadFile(filename)

	if err != nil {
		return err
	}

	return json.Unmarshal(content, view)
}

func (p *Projection) Project(events []interface{}) error {

	state := p.initialState()
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

	err = os.MkdirAll(path.Join(p.root), os.ModePerm)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(path.Join(p.root, p.name+".json"), viewBytes, 0666)
	if err != nil {
		return err
	}

	return nil
}
