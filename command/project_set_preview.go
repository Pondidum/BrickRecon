package command

import (
	"brickrecon/config"
	"brickrecon/domain"
	"brickrecon/lego"
	"brickrecon/storage"
	"brickrecon/tracing"
	"context"
	"fmt"

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

	details   bool
	remaining bool
}

func (c *ProjectSetPreviewCommand) Name() string {
	return "project set preview"
}

func (c *ProjectSetPreviewCommand) Synopsis() string {
	return "shows what adding a set would do to a project"
}

func (c *ProjectSetPreviewCommand) Flags() *pflag.FlagSet {
	flags := pflag.NewFlagSet("project set preview", pflag.ContinueOnError)
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
	// partKey := fmt.Sprintf("%s|%s", args[1], args[2])

	project, err := storage.GetProjectView(ctx, store, storage.WithName(name))
	if err != nil {
		return tracing.Error(span, err)
	}

	set, err := storage.GetLegoSet(ctx, store, setNumber)
	if err != nil {
		return tracing.Error(span, err)
	}

	for _, part := range set.Parts {

		if inv, found := project.Parts[part.Key()]; found {

			currentStock := domain.GetStock(project.Stock, inv.Number, inv.Color)
			if currentStock >= inv.Wanted {
				continue
			}
			newStock := min(currentStock+part.Quantity, inv.Wanted)

			trailer := ""
			if newStock == inv.Wanted {
				trailer = "(complete)"
			}

			fmt.Printf("%s %s: %d/%d (+%d) %s\n",
				lego.GetColorName(part.ColorId), part.Name,
				newStock, inv.Wanted,
				part.Quantity,
				trailer,
			)

		}

	}

	return nil
}
