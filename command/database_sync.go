package command

import (
	"brickrecon/config"
	"brickrecon/storage"
	"brickrecon/tracing"
	"compress/gzip"
	"context"
	"database/sql"
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/spf13/pflag"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

// table dbcolumn csvcolumn
var ColumnSources = map[string]map[string]string{
	"colors": map[string]string{
		"transparent": "is_trans",
	},
	"part_categories": map[string]string{},
	"parts": map[string]string{
		"category_id": "part_cat_id",
	},
	"inventories": map[string]string{},
	"sets":        map[string]string{},
	"inventory_parts": map[string]string{
		"spare":     "is_spare",
		"image_url": "img_url",
	},
}

func NewDatabaseSyncCommand() *DatabaseSyncCommand {
	return &DatabaseSyncCommand{
		tr: otel.Tracer("command.database.sync"),
	}
}

type DatabaseSyncCommand struct {
	tr trace.Tracer
}

func (c *DatabaseSyncCommand) Name() string {
	return "database sync"
}

func (c *DatabaseSyncCommand) Synopsis() string {
	return "sync the rebrickable db locally"
}

func (c *DatabaseSyncCommand) Flags() *pflag.FlagSet {
	flags := pflag.NewFlagSet("project view", pflag.ContinueOnError)
	return flags
}

func (c *DatabaseSyncCommand) Execute(ctx context.Context, config *config.Config, args []string) error {
	ctx, span := c.tr.Start(ctx, "execute")
	defer span.End()

	db, err := storage.NewClient(ctx, config.DatabaseFile)
	if err != nil {
		return tracing.Error(span, err)
	}

	tx, err := db.BeginTx(ctx)
	if err != nil {
		return tracing.Error(span, err)
	}
	defer tx.Rollback()

	if err := c.ensureTables(ctx, tx); err != nil {
		return tracing.Error(span, err)
	}

	fmt.Println("Tables created")

	now := time.Now().UnixMilli()

	for tableName := range ColumnSources {

		fmt.Println("Downloading", tableName, "...")
		res, err := http.Get(fmt.Sprintf("https://cdn.rebrickable.com/media/downloads/%s.csv.gz?%d", tableName, now))
		if err != nil {
			return tracing.Error(span, err)
		}
		defer res.Body.Close()

		gzr, err := gzip.NewReader(res.Body)
		if err != nil {
			return tracing.Error(span, err)
		}

		fmt.Println("Fetched", tableName)

		columnMapping, err := c.createColumnMapping(ctx, tx, tableName)
		if err != nil {
			return tracing.Error(span, err)
		}

		csvr := csv.NewReader(gzr)

		header, err := csvr.Read()
		if err != nil {
			return tracing.Error(span, err)
		}

		stmt := c.insertStatement(tableName, header, columnMapping)
		insert, err := tx.PrepareContext(ctx, stmt)
		if err != nil {
			return tracing.Error(span, err)
		}

		count := 0

		for {
			record, err := csvr.Read()
			if err == io.EOF {
				break
			} else if err != nil {
				return tracing.Error(span, err)
			}

			args := make([]any, 0, len(columnMapping))
			for i, name := range header {
				if _, found := columnMapping[name]; found {
					args = append(args, record[i])
				}
			}

			if _, err := insert.ExecContext(ctx, args...); err != nil {
				return tracing.Error(span, err)
			}

			count++
		}

		if err := insert.Close(); err != nil {
			return tracing.Error(span, err)
		}

		fmt.Println("Inserted", count, tableName)
	}

	fmt.Println("Committing transaction...")
	if err := tx.Commit(); err != nil {
		return tracing.Error(span, err)
	}

	fmt.Println("Compressing db...")
	if err := db.Vacuum(ctx); err != nil {
		return tracing.Error(span, err)
	}

	fmt.Println("Done.")

	return nil
}

func (c *DatabaseSyncCommand) createColumnMapping(ctx context.Context, tx *sql.Tx, tableName string) (map[string]string, error) {

	row, err := tx.QueryContext(ctx,
		"select name from pragma_table_info(@table)",
		sql.Named("table", "rebrickable_"+tableName))
	if err != nil {
		return nil, err
	}
	defer row.Close()

	mapping := map[string]string{}

	for row.Next() {
		dbColumn := ""

		if err := row.Scan(&dbColumn); err != nil {
			return nil, err
		}

		if csvColumn, found := ColumnSources[tableName][dbColumn]; found {
			mapping[csvColumn] = dbColumn
		} else {
			mapping[dbColumn] = dbColumn
		}
	}

	return mapping, nil
}

func (c *DatabaseSyncCommand) insertStatement(tableName string, csvHeaders []string, columnMapping map[string]string) string {

	sb := strings.Builder{}
	sb.WriteString("insert into rebrickable_")
	sb.WriteString(tableName)
	sb.WriteString("(")

	names := make([]string, 0, len(columnMapping))
	for _, h := range csvHeaders {
		if colName, found := columnMapping[h]; found {
			names = append(names, colName)
		}
	}

	sb.WriteString(strings.ToLower(strings.Join(names, ", ")))

	sb.WriteString(")\n")
	sb.WriteString("values(")
	sb.WriteString("?")
	for range len(columnMapping) - 1 {
		sb.WriteString(", ?")
	}
	sb.WriteString(")")

	return sb.String()
}

func (c *DatabaseSyncCommand) ensureTables(ctx context.Context, tx *sql.Tx) error {
	ctx, span := c.tr.Start(ctx, "ensure_tables")
	defer span.End()

	createTables := `
		create table if not exists rebrickable_part_categories (
			id int primary key,
			name text not null
		);

		create table if not exists rebrickable_colors (
			id int primary key,
			name text not null,
			rgb text,
			transparent int
		);
		
		create table if not exists rebrickable_parts (
			part_num text primary key,
			name text not null,
			category_id int references rebrickable_part_categories(id) deferrable initially deferred
		);

		create table if not exists rebrickable_sets (
			set_num text primary key,
			name text not null,
			year int not null,
			theme_id int not null,
			num_parts int not null
		);

		create table if not exists rebrickable_inventories (
			id int primary key,
			version int not null,
			set_num text not null -- note, no fk here as this could also be from the minifigs table
		);

		create table if not exists rebrickable_inventory_parts (
			inventory_id int references rebrickable_inventories(id) deferrable initially deferred,
			part_num text not null references rebrickable_parts(part_num) deferrable initially deferred,
			color_id int not null references rebrickable_colors(id) deferrable initially deferred,
			quantity int not null,
			spare int not null,
			image_url text
		);
	`

	if _, err := tx.ExecContext(ctx, createTables); err != nil {
		return tracing.Error(span, err)
	}

	clearTables := `
		delete from rebrickable_sets where 1=1;
		delete from rebrickable_inventory_parts where 1=1;
		delete from rebrickable_inventories where 1=1;
		delete from rebrickable_parts where 1=1;
		delete from rebrickable_colors where 1=1;
		delete from rebrickable_part_categories where 1=1;
	`

	if _, err := tx.ExecContext(ctx, clearTables); err != nil {
		return tracing.Error(span, err)
	}

	return nil
}
