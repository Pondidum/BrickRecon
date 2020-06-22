package app

import (
	"brickrecon/adapters"
	"brickrecon/background"
	"brickrecon/lego"
	"context"
	"io"

	"github.com/honeycombio/beeline-go"
)

func CreateProject(ctx context.Context, store *AppStore, modelName string, partsFile io.Reader) (func(), error) {

	beeline.AddField(ctx, "model_name", modelName)

	parts, err := adapters.ReadPartsList(partsFile)
	if err != nil {
		beeline.AddField(ctx, "read_parts_error", err)
		return nil, err
	}

	beeline.AddField(ctx, "parts_count", len(parts))

	project := lego.NewProject(modelName, parts)

	if err := store.Save(ctx, project); err != nil {
		beeline.AddField(ctx, "save_project_error", err)
		return nil, err
	}

	wait := store.SendMessage(ctx, &background.PartsAddedMessage{
		Parts: parts,
	})

	return wait, nil
}
