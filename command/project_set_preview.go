package command

import (
	"brickrecon/config"
	"brickrecon/domain"
	"brickrecon/lego"
	"brickrecon/storage"
	"brickrecon/tracing"
	"context"
	"fmt"
	"os"

	"github.com/spf13/pflag"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

func NewProjectSetPreviewCommand() *ProjectSetPreviewCommand {
	return &ProjectSetPreviewCommand{
		tr: otel.Tracer("command.project.set.preview"),
	}
}

type ProjectSetPreviewCommand struct {
	tr trace.Tracer

	renderer string
}

func (c *ProjectSetPreviewCommand) Name() string {
	return "project set preview"
}

func (c *ProjectSetPreviewCommand) Synopsis() string {
	return "shows what adding a set would do to a project"
}

func (c *ProjectSetPreviewCommand) Flags() *pflag.FlagSet {
	flags := pflag.NewFlagSet("project set preview", pflag.ContinueOnError)
	flags.StringVar(&c.renderer, "format", "table", "")
	return flags
}

func (c *ProjectSetPreviewCommand) Execute(ctx context.Context, config *config.Config, args []string) error {
	ctx, span := c.tr.Start(ctx, "execute")
	defer span.End()

	if len(args) != 2 {
		return tracing.Errorf(span, "this command takes exactly 2 arguments: project_name, set_number")
	}

	store, err := storage.NewClient(ctx, config.DatabaseFile)
	if err != nil {
		return tracing.Error(span, err)
	}

	name := args[0]
	setNumber := lego.SetNumber(args[1])

	project, err := storage.GetProjectView(ctx, store, storage.WithName(name))
	if err != nil {
		return tracing.Error(span, err)
	}

	set, err := storage.GetLegoSet(ctx, store, storage.WithSetNumber(setNumber))
	if err != nil {
		return tracing.Error(span, err)
	}

	parts := make([]partDiff, 0, len(set.Parts))

	for _, part := range set.Parts {

		if inv, found := project.Parts[part.Key()]; found {

			currentStock := domain.GetStock(project.Stock, inv.Number, inv.Color)
			if currentStock >= inv.Wanted {
				continue
			}
			newStock := min(currentStock+part.Quantity, inv.Wanted)

			parts = append(parts, partDiff{
				PartNo:   part.Id,
				PartName: part.Name,
				Colour:   lego.GetColorName(part.ColorId),
				Stock:    fmt.Sprintf("%d/%d", newStock, inv.Wanted),
				Change:   part.Quantity,
			})
		}

	}

	if len(parts) > 1 {
		if err := Render(c.renderer, os.Stdout, parts); err != nil {
			return tracing.Error(span, err)
		}
	}

	return nil
}

type partDiff struct {
	PartNo   lego.PartId
	PartName lego.PartName
	Colour   string
	Stock    string
	Change   int
}
