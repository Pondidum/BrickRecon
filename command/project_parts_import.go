package command

import (
	"brickrecon/bricklink"
	"brickrecon/config"
	"brickrecon/domain"
	"brickrecon/storage"
	"brickrecon/tracing"
	"context"
	"database/sql"
	"fmt"
	"os"

	"github.com/google/uuid"
	"github.com/spf13/pflag"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

func NewProjectPartsImportCommand() *ProjectPartsImportCommand {
	return &ProjectPartsImportCommand{
		tr:             otel.Tracer("command.project.parts.import"),
		brickowlApiKey: os.Getenv("BRICKOWL_API_KEY"),
	}
}

type ProjectPartsImportCommand struct {
	tr trace.Tracer

	brickowlApiKey string
}

func (c *ProjectPartsImportCommand) Name() string {
	return "project parts import"
}

func (c *ProjectPartsImportCommand) Synopsis() string {
	return "import a parts list to a project "
}

func (c *ProjectPartsImportCommand) Flags() *pflag.FlagSet {
	flags := pflag.NewFlagSet("project parts import", pflag.ContinueOnError)
	flags.StringVar(&c.brickowlApiKey, "brickowl-apikey", "", "")
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

	project, err := GetProjectByName(ctx, store, name)
	if err != nil {
		return tracing.Error(span, err)
	}

	fmt.Println(project.Name, "currently has", len(project.Parts), "parts")

	content, err := os.Open(wantedList)
	if err != nil {
		return tracing.Error(span, err)
	}
	defer content.Close()

	parts, stock, err := bricklink.ParseWantedList(ctx, content)
	if err != nil {
		return tracing.Error(span, err)
	}

	fmt.Println("Wanted List has", len(parts))

	if err := project.AddParts(parts); err != nil {
		return tracing.Error(span, err)
	}

	if err := project.AddStock(stock); err != nil {
		return tracing.Error(span, err)
	}

	if err := store.SaveAggregate(ctx, project); err != nil {
		return tracing.Error(span, err)
	}

	fmt.Println("Done")

	return nil
}

func GetProjectByName(ctx context.Context, client *storage.Client, name string) (*domain.Project, error) {

	tx, err := client.BeginTx(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	row := tx.QueryRowContext(ctx,
		`select aggregate_id from auto_projections where aggregate_type = 'Project' and view ->> '$.Name' == @name`,
		sql.Named("name", name))

	var id uuid.UUID
	if err := row.Scan(&id); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	project := domain.BlankProject()

	if err := client.LoadAggregate(ctx, id, project); err != nil {
		return nil, err
	}

	return project, nil
}
