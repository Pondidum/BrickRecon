package command

import (
	"brickrecon/config"
	"brickrecon/domain"
	"brickrecon/storage"
	"brickrecon/tracing"
	"brickrecon/util"
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/pflag"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

func NewProjectListCommand() *ProjectListCommand {
	return &ProjectListCommand{
		tr:             otel.Tracer("command.project.list"),
		brickowlApiKey: os.Getenv("BRICKOWL_API_KEY"),
	}
}

type ProjectListCommand struct {
	tr trace.Tracer

	brickowlApiKey string
}

func (c *ProjectListCommand) Name() string {
	return "project list"
}

func (c *ProjectListCommand) Synopsis() string {
	return "list all projects"
}

func (c *ProjectListCommand) Flags() *pflag.FlagSet {
	flags := pflag.NewFlagSet("project list", pflag.ContinueOnError)
	flags.StringVar(&c.brickowlApiKey, "brickowl-apikey", "", "")
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

	projects, err := GetAllProjects(ctx, store)
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

func GetAllProjects(ctx context.Context, client *storage.Client) ([]*domain.ProjectView, error) {

	tx, err := client.BeginTx(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	row, err := tx.QueryContext(ctx, `select view from auto_projections where aggregate_type = 'Project'`)
	if err != nil {
		return nil, err
	}

	projects := []*domain.ProjectView{}
	for row.Next() {

		var viewJson []byte
		if err := row.Scan(&viewJson); err != nil {
			return nil, err
		}

		view := &domain.ProjectView{}
		if err := json.Unmarshal(viewJson, view); err != nil {
			return nil, err
		}

		projects = append(projects, view)
	}

	return projects, nil
}
