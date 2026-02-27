package command

import (
	"brickrecon/config"
	"brickrecon/domain"
	"brickrecon/storage"
	"brickrecon/tracing"
	"context"
	"fmt"

	"github.com/spf13/pflag"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

func NewProjectNewCommand() *ProjectNewCommand {
	return &ProjectNewCommand{
		tr: otel.Tracer("command.project.new"),
	}
}

type ProjectNewCommand struct {
	tr trace.Tracer
}

func (c *ProjectNewCommand) Name() string {
	return "project new"
}

func (c *ProjectNewCommand) Synopsis() string {
	return "create a new project"
}

func (c *ProjectNewCommand) Flags() *pflag.FlagSet {
	flags := pflag.NewFlagSet("project new", pflag.ContinueOnError)
	return flags
}

func (c *ProjectNewCommand) Execute(ctx context.Context, config *config.Config, args []string) error {
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

	view, err := storage.GetProjectView(ctx, store, storage.IncludeArchived(), storage.WithName(name))
	if err != nil && err != storage.ErrViewNotFound {
		return tracing.Error(span, err)
	}
	if view != nil {
		return tracing.Errorf(span, "There is already a project called %s", name)
	}

	project, err := domain.CreateProject(name)
	if err != nil {
		return tracing.Error(span, err)
	}

	if err := store.SaveAggregate(ctx, project); err != nil {
		return tracing.Error(span, err)
	}

	fmt.Println("Done")

	return nil
}
