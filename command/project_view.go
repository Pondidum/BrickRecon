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

func NewProjectViewCommand() *ProjectViewCommand {
	return &ProjectViewCommand{
		tr: otel.Tracer("command.project.view"),
	}
}

type ProjectViewCommand struct {
	tr trace.Tracer

	details   bool
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
	flags.BoolVar(&c.details, "details", false, "")
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

	if c.details {
		for _, part := range project.Parts {
			fmt.Println(part.Number, part.Color, part.Wanted, domain.GetStock(project.Stock, part.Number, part.Color))
		}
	} else if c.remaining {
		lines := make([]string, 0, len(project.Parts)+1)

		lines = append(lines, "Part No | Part Name | Colour | Wanted | Owned")

		for _, part := range project.Parts {
			stock := domain.GetStock(project.Stock, part.Number, part.Color)

			if stock >= part.Wanted {
				continue
			}

			color := lego.GetColorName(part.Color)

			lines = append(lines, fmt.Sprintf("%s | %s | %s | %d | %d", part.Number, "?", color, part.Wanted, stock))
		}

		fmt.Println(util.TableOutput(lines))
	}

	return nil
}
