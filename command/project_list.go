package command

import (
	"brickrecon/config"
	"brickrecon/storage"
	"brickrecon/tracing"
	"brickrecon/util"
	"context"
	"fmt"

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
	tr trace.Tracer
}

func (c *ProjectListCommand) Name() string {
	return "project list"
}

func (c *ProjectListCommand) Synopsis() string {
	return "list all projects"
}

func (c *ProjectListCommand) Flags() *pflag.FlagSet {
	flags := pflag.NewFlagSet("project list", pflag.ContinueOnError)
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

	projects, err := storage.GetProjectViewsAll(ctx, store)
	if err != nil {
		return tracing.Error(span, err)
	}

	lines := make([]string, len(projects)+1)
	lines[0] = "Name | Unique Parts | Total Parts | Owned Parts"

	for i, project := range projects {
		lines[i+1] = fmt.Sprintf("%s | %d | %d | %d", project.Name, project.UniqueParts(), project.TotalParts(), project.OwnedParts())
	}

	fmt.Println(util.TableOutput(lines))

	return nil
}
