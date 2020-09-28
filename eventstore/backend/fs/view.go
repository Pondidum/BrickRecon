package fs

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"os"
)

type FsView struct {
	filename string
}

func (v *FsView) LastEventIndex(ctx context.Context) (int, error) {
	return readCheckIndex(v.filename)
}

func (v *FsView) ReadView(ctx context.Context, view interface{}) error {

	content, err := ioutil.ReadFile(v.filename)

	if err != nil {

		if os.IsNotExist(err) {
			return nil
		}

		return err
	}

	return json.Unmarshal(content, view)
}

func (v *FsView) WriteView(ctx context.Context, view interface{}, lastIndex int) error {
	viewBytes, err := json.MarshalIndent(view, "", "  ")

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
