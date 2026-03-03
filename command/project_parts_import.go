package command

import (
	"brickrecon/bricklink"
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

func NewProjectPartsImportCommand() *ProjectPartsImportCommand {
	return &ProjectPartsImportCommand{
		tr: otel.Tracer("command.project.parts.import"),
	}
}

type ProjectPartsImportCommand struct {
	tr     trace.Tracer
	dryrun bool
}

func (c *ProjectPartsImportCommand) Name() string {
	return "project parts import"
}

func (c *ProjectPartsImportCommand) Synopsis() string {
	return "import a parts list to a project "
}

func (c *ProjectPartsImportCommand) Flags() *pflag.FlagSet {
	flags := pflag.NewFlagSet("project parts import", pflag.ContinueOnError)
	flags.BoolVar(&c.dryrun, "dry-run", false, "")
	return flags
}

func (c *ProjectPartsImportCommand) Execute(ctx context.Context, config *config.Config, args []string) error {
	ctx, span := c.tr.Start(ctx, "execute")
	defer span.End()

	if len(args) != 2 {
		return tracing.Errorf(span, "this command takes exactly 2 arguments: name and parts-path")
	}

	store, err := storage.NewClient(ctx, config.DatabaseFile)
	if err != nil {
		return tracing.Error(span, err)
	}

	name := args[0]
	wantedList := args[1]

	project, err := storage.GetProjectByName(ctx, store, name)
	if err != nil {
		return tracing.Error(span, err)
	}

	fmt.Println(project.Name, "currently has", len(project.Parts), "parts")

	content, err := os.Open(wantedList)
	if err != nil {
		return tracing.Error(span, err)
	}
	defer content.Close()

	getPart := func(partId lego.PartId) (*lego.Part, error) {
		parts, err := storage.FindMatchingParts(ctx, store, partId)
		if err != nil {
			return nil, err
		}
		switch len(parts) {
		case 0:
			return nil, fmt.Errorf("unable to find any bricks matching %s", partId)
		case 1:
			return parts[0], nil
		default:
			// we could query the user here to pick the right one
			return parts[0], nil
		}
	}

	parts, stock, err := bricklink.ParseWantedList(ctx, getPart, content)
	if err != nil {
		return tracing.Error(span, err)
	}

	fmt.Println("Wanted List has", len(parts))

	if c.dryrun {
		fmt.Println("Would import", len(parts))

		for _, part := range parts {
			fmt.Println(part.Id, part.Name, lego.GetColorName(part.ColorId), part.Quantity, domain.GetStock(stock, part.Id, part.ColorId))
		}

		return nil
	}

	fmt.Printf("Adding parts...")
	if err := project.AddParts(parts); err != nil {
		return tracing.Error(span, err)
	}
	fmt.Println("Done")

	fmt.Print("Adding Stock...")
	if err := project.AddStock(stock); err != nil {
		return tracing.Error(span, err)
	}
	fmt.Println("Done")

	if err := store.SaveAggregate(ctx, project); err != nil {
		return tracing.Error(span, err)
	}

	fmt.Println("Done")

	return nil
}
