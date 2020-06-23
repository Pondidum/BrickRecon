package app

import (
	"brickrecon/adapters"
	"brickrecon/background"
	"brickrecon/lego"
	"context"
	"os"

	"github.com/honeycombio/beeline-go"
)

func ImportKit(ctx context.Context, store *AppStore, kitNumber string) (func(), error) {

	beeline.AddField(ctx, "kit_number", kitNumber)

	api := adapters.NewBrickOwlApi(os.Getenv("BRICKOWL_API_KEY"))

	parts, err := api.GetParts(kitNumber)
	if err != nil {
		beeline.AddField(ctx, "fetch_parts_error", err)
		return nil, err
	}

	name, err := api.GetSetName(kitNumber)
	if err != nil {
		beeline.AddField(ctx, "fetch_set_name_error", err)
		return nil, err
	}

	beeline.AddField(ctx, "kit_name", name)
	beeline.AddField(ctx, "parts_count", len(parts))

	kit := lego.ImportKit(kitNumber, name, parts)

	if err := store.Save(ctx, kit); err != nil {
		beeline.AddField(ctx, "save_kit_error", err)
		return nil, err
	}

	wait := store.SendMessage(ctx, &background.PartsAddedMessage{
		Parts: parts,
	})

	return wait, nil
}
