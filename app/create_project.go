package app

import (
	"brickrecon/lego"
	"brickrecon/stud_io"
	"context"
	"io"

	"github.com/honeycombio/beeline-go"
)

func CreateProject(ctx context.Context, store *AppStore, projectName lego.ProjectName, partsFile io.Reader) (func(), error) {

	beeline.AddField(ctx, "project_name", projectName)

	parts, err := stud_io.ReadPartsList(partsFile)
	if err != nil {
		beeline.AddField(ctx, "read_parts_error", err)
		return nil, err
	}

	beeline.AddField(ctx, "parts_count", len(parts))

	builder, err := NewPartBuilder(ctx, store.EventStore)
	if err != nil {
		beeline.AddField(ctx, "builder_error", err)
		return nil, err
	}
	for _, part := range parts {
		err := builder.StorePart(ctx, part)

		if err != nil {
			beeline.AddField(ctx, string(part.Key)+"_error", err)
		}
	}

	project := lego.NewProject(projectName, parts)

	if err := store.Save(ctx, project); err != nil {
		beeline.AddField(ctx, "save_project_error", err)
		return nil, err
	}

	return func() {}, nil
}
