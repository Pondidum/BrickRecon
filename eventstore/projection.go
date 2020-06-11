package eventstore

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type View interface {
	LastEventIndex() (int, error)
	ReadView(view interface{}) error
	WriteView(view interface{}, lastIndex int) error
}

type FsView struct {
	filename string
}

func (v *FsView) LastEventIndex() (int, error) {
	return readCheckIndex(v.filename)
}

func (v *FsView) ReadView(view interface{}) error {

	content, err := ioutil.ReadFile(v.filename)

	if err != nil {

		if os.IsNotExist(err) {
			return nil
		}

		return err
	}

	return json.Unmarshal(content, view)
}

func (v *FsView) WriteView(view interface{}, lastIndex int) error {
	viewBytes, err := json.Marshal(view)

	if err != nil {
		return err
	}

	if err := ioutil.WriteFile(v.filename, viewBytes, 0666); err != nil {
		return err
	}

	if err := writeCheckIndex(v.filename, lastIndex); err != nil {
		return err
	}

	return nil
}
