package create

import (
	"io"
	"mvc/app"
	"mvc/background"
	"mvc/lego"
)

func CreateProject(store *app.AppStore, modelName string, partsFile io.Reader) (func(), error) {

	parts, err := lego.ReadPartsList(partsFile)
	if err != nil {
		return nil, err
	}

	project := lego.NewProject(modelName, parts)

	if err := store.Save(project); err != nil {
		return nil, err
	}

	wait := store.SendMessage(&background.PartsAddedMessage{
		Parts: parts,
	})

	return wait, nil
}
