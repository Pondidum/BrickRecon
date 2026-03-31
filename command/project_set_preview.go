package command

import (
	"brickrecon/config"
	"brickrecon/domain"
	"brickrecon/lego"
	"brickrecon/storage"
	"brickrecon/tracing"
	"brickrecon/util"
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

	set, err := storage.GetLegoSet(ctx, store, storage.WithSetNumber(setNumber))
	if err != nil {
		return tracing.Error(span, err)
	}

	lines := make([]string, 0, len(set.Parts)+1)
	lines = append(lines, "PartNo | Part Name | Color | Stock | Change")

	for _, part := range set.Parts {

		if inv, found := project.Parts[part.Key()]; found {

			currentStock := domain.GetStock(project.Stock, inv.Number, inv.Color)
			if currentStock >= inv.Wanted {
				continue
			}
			newStock := min(currentStock+part.Quantity, inv.Wanted)

			lines = append(lines, fmt.Sprintf("%s | %s | %s | %d/%d | +%d",
				part.Id,
				part.Name,
				lego.GetColorName(part.ColorId),
				newStock, inv.Wanted,
				part.Quantity,
			))

		}

	}

	if len(lines) > 1 {
		fmt.Println(util.TableOutput(lines))
	}

	return nil
}
