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

func NewProjectViewCommand() *ProjectViewCommand {
	return &ProjectViewCommand{
		tr: otel.Tracer("command.project.view"),
	}
}

type ProjectViewCommand struct {
	tr trace.Tracer

	renderer  string
	remaining bool
}

func (c *ProjectViewCommand) Name() string {
	return "project view"
}

func (c *ProjectViewCommand) Synopsis() string {
	return "view a project"
}

func (c *ProjectViewCommand) Flags() *pflag.FlagSet {
	flags := pflag.NewFlagSet("project view", pflag.ContinueOnError)
	flags.StringVar(&c.renderer, "format", "table", "")
	flags.BoolVar(&c.remaining, "remaining", false, "")
	return flags
}

func (c *ProjectViewCommand) Execute(ctx context.Context, config *config.Config, args []string) error {
	ctx, span := c.tr.Start(ctx, "execute")
	defer span.End()

	if len(args) != 1 {
		return tracing.Errorf(span, "this command takes exactly 1 argument: name")
	}

	store, err := storage.NewClient(ctx, config.DatabaseFile)
	if err != nil {
		return tracing.Error(span, err)
	}

	name := args[0]

	project, err := storage.GetProjectView(ctx, store, storage.WithName(name))
	if err != nil {
		return tracing.Error(span, err)
	}

	fmt.Println(project.UniqueParts(), " unique parts, ", project.TotalParts(), " total parts, ", project.OwnedParts(), " stocked parts")

	type PartSummary struct {
		PartNo   lego.PartId
		PartName lego.PartName
		Colour   string
		Wanted   int
		Stock    int
	}

	summary := make([]PartSummary, 0, len(project.Parts))
	for _, part := range project.Parts {
		stock := domain.GetStock(project.Stock, part.Number, part.Color)

		if c.remaining && stock >= part.Wanted {
			continue
		}

		color := lego.GetColorName(part.Color)

		summary = append(summary, PartSummary{
			PartNo:   part.Number,
			PartName: "?",
			Colour:   color,
			Wanted:   part.Wanted,
			Stock:    stock,
		})
	}

	if err := Render(c.renderer, os.Stdout, summary); err != nil {
		return tracing.Error(span, err)
	}

	return nil
}
