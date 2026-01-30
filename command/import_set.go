package command

import (
	"brickrecon/brickowl"
	"brickrecon/config"
	"brickrecon/lego"
	"brickrecon/storage"
	"brickrecon/tracing"
	"context"
	"database/sql"
	"fmt"
	"os"

	"github.com/spf13/pflag"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

func NewImportSetCommand() *ImportSetCommand {
	return &ImportSetCommand{
		tr:             otel.Tracer("command.import.set"),
		brickowlApiKey: os.Getenv("BRICKOWL_API_KEY"),
	}
}

type ImportSetCommand struct {
	tr trace.Tracer

	brickowlApiKey string
}

func (c *ImportSetCommand) Name() string {
	return "import set"
}

func (c *ImportSetCommand) Synopsis() string {
	return "import a lego set"
}

func (c *ImportSetCommand) Flags() *pflag.FlagSet {
	flags := pflag.NewFlagSet("import set", pflag.ContinueOnError)
	flags.StringVar(&c.brickowlApiKey, "brickowl-apikey", "", "")
	return flags
}

func (c *ImportSetCommand) Execute(ctx context.Context, config *config.Config, args []string) error {
	ctx, span := c.tr.Start(ctx, "execute")
	defer span.End()

	if len(args) != 1 {
		return tracing.Errorf(span, "this command takes exactly 1 argument: set_id")
	}

	apikey := c.brickowlApiKey
	if apikey == "" {
		if val := os.Getenv("BRICKOWL_API_KEY"); val == "" {
			return tracing.Errorf(span, "this command requires a brickowl apikey")
		} else {
			apikey = val
		}
	}

	setNumber := lego.SetNumber(args[0])
	owl := brickowl.NewBrickOwlApi(apikey)

	legoSet, err := owl.GetSet(ctx, setNumber)
	if err != nil {
		return tracing.Error(span, err)
	}

	fmt.Println("SetName:", legoSet.Name)
	fmt.Println("Parts:", len(legoSet.Parts))

	store, err := storage.NewClient(ctx, config.DatabaseFile)
	if err != nil {
		return tracing.Error(span, err)
	}

	tx, err := store.BeginTx(ctx)
	if err != nil {
		return tracing.Error(span, err)
	}
	defer tx.Rollback()

	if _, err := storage.GetLegoSetByNumber(ctx, tx, setNumber); err != nil && err != sql.ErrNoRows {
		fmt.Println("Set is already in the catalogue")
		return nil
	}

	if err := storage.InsertLegoSet(ctx, tx, legoSet); err != nil {
		return tracing.Error(span, err)
	}

	if err := tx.Commit(); err != nil {
		return tracing.Error(span, err)
	}

	fmt.Println("Done")

	return nil
}
