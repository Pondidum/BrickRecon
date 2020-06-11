package eventstore

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
)

type Projector func(state interface{}, event Event) interface{}

type Projection interface {
	ReadView(view interface{}) error
	Project(events []Event) error
	LastEventIndex() (int, error)
}

type FsProjection struct {
	path            string
	initialiseState Initialiser
	projector       Projector
}

func NewProjection(root string, name string, initialiseState Initialiser, projector Projector) Projection {
	filepath := path.Join(root, name+".json")

	return &FsProjection{
		path:            filepath,
		initialiseState: initialiseState,
		projector:       projector,
	}
}

func (p *FsProjection) LastEventIndex() (int, error) {
	return readCheckIndex(p.path)
}

func (p *FsProjection) ReadView(view interface{}) error {

	content, err := ioutil.ReadFile(p.path)

	if err != nil {
		return err
	}

	return json.Unmarshal(content, view)
}

func (p *FsProjection) Project(events []Event) error {

	if len(events) == 0 {
		return nil
	}

	state := p.initialiseState()
	err := p.ReadView(state)

	if err != nil && !os.IsNotExist(err) {
		return err
	}

	var lastIndex int
	for _, e := range events {
		state = p.projector(state, e)
		lastIndex = e.event().Version
	}

	viewBytes, err := json.Marshal(state)

	if err != nil {
		return err
	}

	if err := ioutil.WriteFile(p.path, viewBytes, 0666); err != nil {
		return err
	}

	if err := writeCheckIndex(p.path, lastIndex); err != nil {
		return err
	}

	return nil
}
