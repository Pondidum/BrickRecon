package command

import (
	"brickrecon/config"
	"brickrecon/ldraw"
	"brickrecon/storage"
	"brickrecon/tracing"
	"compress/gzip"
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/spf13/pflag"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

func NewDatabaseSyncCommand() *DatabaseSyncCommand {
	return &DatabaseSyncCommand{
		tr: otel.Tracer("command.database.sync"),
	}
}

type DatabaseSyncCommand struct {
	tr     trace.Tracer
	dryrun bool
	debug  bool
	wipe   bool
}

func (c *DatabaseSyncCommand) Name() string {
	return "database sync"
}

func (c *DatabaseSyncCommand) Synopsis() string {
	return "sync the rebrickable db locally"
}

func (c *DatabaseSyncCommand) Flags() *pflag.FlagSet {
	flags := pflag.NewFlagSet("project view", pflag.ContinueOnError)
	flags.BoolVar(&c.dryrun, "dry-run", false, "")
	flags.BoolVar(&c.debug, "debug", false, "")
	flags.BoolVar(&c.wipe, "wipe", false, "")
	return flags
}

func (c *DatabaseSyncCommand) Execute(ctx context.Context, config *config.Config, args []string) error {
	ctx, span := c.tr.Start(ctx, "execute")
	defer span.End()

	if len(args) != 1 || (args[0] != "ldraw" && args[0] != "rebrickable") {
		return fmt.Errorf("you must specify 'ldraw' or 'rebrickable'")
	}

	source := args[0]

	db, err := storage.NewClient(ctx, config.DatabaseFile)
	if err != nil {
		return tracing.Error(span, err)
	}

	var writer writer

	if c.dryrun {
		writer = &drywriter{debug: c.debug}
	} else {
		writer, err = NewDbWriter(ctx, db)
		if err != nil {
			return tracing.Error(span, err)
		}
	}

	defer writer.Cancel()

	if err := writer.CreateTables(ctx, source); err != nil {
		return tracing.Error(span, err)
	}

	if err := writer.ClearTables(ctx, source); err != nil {
		return tracing.Error(span, err)
	}

	if !c.wipe {

		switch source {
		case "ldraw":
			if err := c.syncLDraw(ctx, writer); err != nil {
				return tracing.Error(span, err)
			}

		case "rebrickable":

			for tableName := range Tables {
				if !strings.HasPrefix(tableName, "rebrickable_") {
					continue
				}
				if err := c.syncTable(ctx, writer, tableName); err != nil {
					return tracing.Error(span, err)
				}
			}

		default:
			return tracing.Errorf(span, "unsupported source %s", source)

		}
	}
	if err := writer.Finish(ctx); err != nil {
		return tracing.Error(span, err)
	}

	fmt.Println("Done.")

	return nil
}

func (c *DatabaseSyncCommand) syncTable(ctx context.Context, writer writer, tableName string) error {
	ctx, span := c.tr.Start(ctx, "sync_table")
	defer span.End()

	filename := strings.TrimPrefix(tableName, "rebrickable_")

	fmt.Print("Downloading ", filename, "...")
	now := time.Now().UnixMilli()
	res, err := http.Get(fmt.Sprintf("https://cdn.rebrickable.com/media/downloads/%s.csv.gz?%d", filename, now))
	if err != nil {
		return tracing.Error(span, err)
	}
	defer res.Body.Close()

	gzr, err := gzip.NewReader(res.Body)
	if err != nil {
		return tracing.Error(span, err)
	}

	fmt.Println("Done")

	csvr := csv.NewReader(gzr)

	header, err := csvr.Read()
	if err != nil {
		return tracing.Error(span, err)
	}

	insert, done, err := writer.Prepare(ctx, tableName, header)
	if err != nil {
		return tracing.Error(span, err)
	}
	defer done()

	for {
		record, err := csvr.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return tracing.Error(span, err)
		}

		if err := insert(record); err != nil {
			return tracing.Error(span, err)
		}
	}

	return nil
}

func (c *DatabaseSyncCommand) syncLDraw(ctx context.Context, writer writer) error {
	ctx, span := c.tr.Start(ctx, "sync_ldraw")
	defer span.End()

	var content io.ReadCloser
	file, err := os.Open("ldraw/complete.zip")
	if err == nil {
		fmt.Println("Using local complete.zip file")
		content = file
	} else {
		fmt.Println("Using remote complete.zip file")
		res, err := http.Get("https://library.ldraw.org/library/updates/complete.zip")
		if err != nil {
			return tracing.Error(span, err)
		}
		content = res.Body
	}
	defer content.Close()

	parts, err := ldraw.ParseDatabaseArchive(ctx, content)
	if err != nil {
		return tracing.Error(span, err)
	}

	insertMoves, movesDone, err := writer.Prepare(ctx, "ldraw_moved_parts", []string{"old_part_num", "new_part_num"})
	if err != nil {
		return tracing.Error(span, err)
	}
	defer movesDone()

	insertAlternates, alternatesDone, err := writer.Prepare(ctx, "ldraw_alternate_ids", []string{"part_num", "alt_part_num"})
	if err != nil {
		return tracing.Error(span, err)
	}
	defer alternatesDone()

	for old, newPart := range parts {
		if newPart.MovedTo != "" {

			if err := insertMoves([]string{old, newPart.MovedTo}); err != nil {
				return tracing.Error(span, err)
			}
		}

		for _, alternate := range newPart.AlternateIds {
			if err := insertAlternates([]string{old, alternate}); err != nil {
				return tracing.Error(span, err)
			}
		}
	}

	return nil
}

type writer interface {
	CreateTables(ctx context.Context, source string) error
	ClearTables(ctx context.Context, source string) error
	Cancel() error
	Prepare(ctx context.Context, tableName string, header []string) (func(record []string) error, func() error, error)
	Finish(ctx context.Context) error
}
