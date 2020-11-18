package app

import (
	"brickrecon/brickowl"
	"brickrecon/lego"
	"context"
	"os"

	"github.com/honeycombio/beeline-go"
)

func ImportKit(ctx context.Context, store *AppStore, kitNumber lego.KitNumber) (func(), error) {

	beeline.AddField(ctx, "kit_number", kitNumber)

	api := brickowl.NewBrickOwlApi(os.Getenv("BRICKOWL_API_KEY"))

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

	builder, err := NewPartBuilder(ctx, store.EventStore)
	if err != nil {
		beeline.AddField(ctx, "builder_error", err)
		return nil, err
	}

	keys := map[lego.PartKey]int{}

	for _, owlPart := range parts {
		builder.FromBrickOwl(ctx, owlPart)
		keys[owlPart.Key] = owlPart.Quantity
	}

	kit := lego.ImportKit(kitNumber, name, keys)

	if err := store.Save(ctx, kit); err != nil {
		beeline.AddField(ctx, "save_kit_error", err)
		return nil, err
	}

	return func() {}, nil
}
