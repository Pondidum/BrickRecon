package command

import (
	"brickrecon/config"
	"brickrecon/storage"
	"brickrecon/tracing"
	"context"
	"fmt"

	"github.com/spf13/pflag"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

func NewProjectRmCommand() *ProjectRmCommand {
	return &ProjectRmCommand{
		tr: otel.Tracer("command.project.rm"),
	}
}

type ProjectRmCommand struct {
	tr         trace.Tracer
	hardDelete bool
}

func (c *ProjectRmCommand) Name() string {
	return "project rm"
}

func (c *ProjectRmCommand) Synopsis() string {
	return "delete a project"
}

func (c *ProjectRmCommand) Flags() *pflag.FlagSet {
	flags := pflag.NewFlagSet("project rm", pflag.ContinueOnError)
	flags.BoolVar(&c.hardDelete, "hard", false, "remove the project entirely, rather than softdelete")
	return flags
}

func (c *ProjectRmCommand) Execute(ctx context.Context, config *config.Config, args []string) error {
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

	project, err := storage.GetProjectByName(ctx, store, name)
	if err != nil {
		return tracing.Error(span, err)
	}

	if c.hardDelete {
		if err := store.DeleteAggregate(ctx, project); err != nil {
			return tracing.Error(span, err)
		}

		return nil
	} else {

		if err := project.Archive(); err != nil {
			return tracing.Error(span, err)
		}

		if err := store.SaveAggregate(ctx, project); err != nil {
			return tracing.Error(span, err)
		}
	}

	fmt.Println("Done")

	return nil
}
