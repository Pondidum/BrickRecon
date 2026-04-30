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

func NewProjectSetAddCommand() *ProjectSetAddCommand {
	return &ProjectSetAddCommand{
		tr: otel.Tracer("command.project.set.preview"),
	}
}

type ProjectSetAddCommand struct {
	tr trace.Tracer

	preview  bool
	renderer string
}

func (c *ProjectSetAddCommand) Name() string {
	return "project set add"
}

func (c *ProjectSetAddCommand) Synopsis() string {
	return "addds the stock needed from a set"
}

func (c *ProjectSetAddCommand) Flags() *pflag.FlagSet {
	flags := pflag.NewFlagSet("project set preview", pflag.ContinueOnError)
	flags.StringVar(&c.renderer, "format", "table", "")
	flags.BoolVar(&c.preview, "preview", false, "show what would be added")
	return flags
}

func (c *ProjectSetAddCommand) Execute(ctx context.Context, config *config.Config, args []string) error {
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

	project, err := storage.GetProjectByName(ctx, store, name)
	if err != nil {
		return tracing.Error(span, err)
	}

	set, err := storage.GetLegoSet(ctx, store, storage.WithSetNumber(setNumber))
	if err != nil {
		return tracing.Error(span, err)
	}

	fmt.Printf("Set %s: %s contains %d unique parts\n", set.Number, set.Name, len(set.Parts))

	toAdd := domain.Stock{}

	for _, part := range set.Parts {

		if inv, found := project.Parts[part.Key()]; found {

			currentStock := domain.GetStock(project.Stock, inv.Number, inv.Color)
			if currentStock >= inv.Wanted {
				continue
			}

			domain.AddStock(toAdd, part.Id, part.ColorId, part.Quantity)
		}
	}

	if len(toAdd) == 0 {
		fmt.Println("No stock to add to the project")
		return nil
	}

	if c.preview {

		diff := make([]partDiff, 0, len(toAdd))

		for partId, colors := range toAdd {
			for colorId, quantity := range colors {

				inv := project.Parts[fmt.Sprintf("%s|%s", partId, colorId)]
				currentStock := domain.GetStock(project.Stock, inv.Number, inv.Color)

				diff = append(diff, partDiff{
					PartNo: partId,
					// PartName: inv.,c
					Colour: lego.GetColorName(colorId),
					Stock:  fmt.Sprintf("%d/%d", currentStock+quantity, inv.Wanted),
					Change: quantity,
				})

			}
		}

		if err := Render(c.renderer, os.Stdout, diff); err != nil {
			return tracing.Error(span, err)
		}

		return nil
	}

	return nil
}
