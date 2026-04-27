package command

import (
	"brickrecon/config"
	"brickrecon/storage"
	"brickrecon/tracing"
	"context"
	"os"

	"github.com/spf13/pflag"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

func NewProjectListCommand() *ProjectListCommand {
	return &ProjectListCommand{
		tr: otel.Tracer("command.project.list"),
	}
}

type ProjectListCommand struct {
	tr       trace.Tracer
	renderer string
	archived bool
}

func (c *ProjectListCommand) Name() string {
	return "project list"
}

func (c *ProjectListCommand) Synopsis() string {
	return "list all projects"
}

func (c *ProjectListCommand) Flags() *pflag.FlagSet {
	flags := pflag.NewFlagSet("project list", pflag.ContinueOnError)
	flags.StringVar(&c.renderer, "format", "table", "")
	flags.BoolVar(&c.archived, "archived", false, "include archived projects")
	return flags
}

func (c *ProjectListCommand) Execute(ctx context.Context, config *config.Config, args []string) error {
	ctx, span := c.tr.Start(ctx, "execute")
	defer span.End()

	if len(args) != 0 {
		return tracing.Errorf(span, "this command takes no arguments")
	}

	store, err := storage.NewClient(ctx, config.DatabaseFile)
	if err != nil {
		return tracing.Error(span, err)
	}

	opts := []storage.QueryOption{}
	if c.archived {
		opts = append(opts, storage.IncludeArchived())
	}

	views, err := storage.GetProjectViews(ctx, store, opts...)
	if err != nil {
		return tracing.Error(span, err)
	}

	projects := make([]projectSummary, 0, len(views))

	for _, view := range views {
		projects = append(projects, projectSummary{
			Name:        view.Name,
			UniqueParts: view.UniqueParts(),
			TotalParts:  view.TotalParts(),
			Stock:       view.OwnedParts(),
		})
	}

	if err := Render(c.renderer, os.Stdout, projects); err != nil {
		return tracing.Error(span, err)
	}

	return nil
}

type projectSummary struct {
	Name        string
	UniqueParts int
	TotalParts  int
	Stock       int
}
