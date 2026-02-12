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
	"os"
	"strings"
	"time"

	"github.com/spf13/pflag"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

func NewDatabaseSyncCommand() *DatabaseSyncCommand {
	return &DatabaseSyncCommand{
		tr:                otel.Tracer("command.database.sync"),
		rebrickableApiKey: os.Getenv("REBRICKABLE_API_KEY"),
	}
}

type DatabaseSyncCommand struct {
	tr trace.Tracer

	rebrickableApiKey string
}

func (c *DatabaseSyncCommand) Name() string {
	return "database sync"
}

func (c *DatabaseSyncCommand) Synopsis() string {
	return "sync the rebrickable db locally"
}

func (c *DatabaseSyncCommand) Flags() *pflag.FlagSet {
	flags := pflag.NewFlagSet("project view", pflag.ContinueOnError)
	flags.StringVar(&c.rebrickableApiKey, "rebrickable-apikey", "", "")
	return flags
}

func (c *DatabaseSyncCommand) Execute(ctx context.Context, config *config.Config, args []string) error {
	ctx, span := c.tr.Start(ctx, "execute")
	defer span.End()

	db, err := storage.NewClient(ctx, config.DatabaseFile)
	if err != nil {
		return tracing.Error(span, err)
	}

	tx, _ := db.BeginTx(ctx)
	if err != nil {
		return tracing.Error(span, err)
	}
	defer tx.Rollback()

	if err := c.ensureTables(ctx, tx); err != nil {
		return tracing.Error(span, err)
	}

	fmt.Println("Tables created")

	mappings := map[string]map[string]string{
		"colors": map[string]string{
			"transparent": "is_trans",
		},
		"part_categories": map[string]string{},
		"parts": map[string]string{
			"category_id": "part_cat_id",
		},
		"inventories": map[string]string{},
		"inventory_parts": map[string]string{
			"spare":     "is_spare",
			"image_url": "img_url",
		},
	}

	files := []string{"colors", "part_categories", "parts", "inventories", "inventory_parts"}
	now := time.Now().UnixMilli()

	for _, file := range files {

		fmt.Println("Downloading", file, "...")
		res, err := http.Get(fmt.Sprintf("https://cdn.rebrickable.com/media/downloads/%s.csv.gz?%d", file, now))
		if err != nil {
			return tracing.Error(span, err)
		}
		defer res.Body.Close()

		gzr, err := gzip.NewReader(res.Body)
		if err != nil {
			return tracing.Error(span, err)
		}

		cols := map[string]string{}

		row, err := tx.QueryContext(ctx, "select name from pragma_table_info(@table)", sql.Named("table", "rebrickable_"+file))
		if err != nil {
			return tracing.Error(span, err)
		}
		defer row.Close()
		for row.Next() {
			column := ""
			if err := row.Scan(&column); err != nil {
				return tracing.Error(span, err)
			}

			if name, found := mappings[file][column]; found {
				cols[name] = column
			} else {
				cols[column] = column
			}
		}

		csvr := csv.NewReader(gzr)

		header, err := csvr.Read()
		if err != nil {
			return tracing.Error(span, err)
		}

		sb := strings.Builder{}
		sb.WriteString("insert into rebrickable_")
		sb.WriteString(file)
		sb.WriteString("(")

		names := make([]string, 0, len(cols))
		for _, h := range header {
			if colName, found := cols[h]; found {
				names = append(names, colName)
			}
		}

		sb.WriteString(strings.ToLower(strings.Join(names, ", ")))

		sb.WriteString(")\n")
		sb.WriteString("values(")
		sb.WriteString("?")
		for range len(cols) - 1 {
			sb.WriteString(", ?")
		}
		sb.WriteString(")")

		insert, err := tx.PrepareContext(ctx, sb.String())
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

			args := make([]any, 0, len(cols))
			for i, name := range header {
				if _, found := cols[name]; found {
					args = append(args, record[i])
				}
			}

			if _, err := insert.ExecContext(ctx, args...); err != nil {
				return tracing.Error(span, err)
			}

			count++
		}

		fmt.Println("Inserted", count, file)
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

func (c *DatabaseSyncCommand) ensureTables(ctx context.Context, tx *sql.Tx) error {
	ctx, span := c.tr.Start(ctx, "ensure_tables")
	defer span.End()

	// dropTables := `
	// 	drop table if exists rebrickable_inventories;
	// 	drop table if exists rebrickable_inventory_parts;
	// 	drop table if exists rebrickable_parts;
	// 	drop table if exists rebrickable_part_categories;
	// 	drop table if exists rebrickable_colors;
	// `

	// if _, err := tx.ExecContext(ctx, dropTables); err != nil {
	// 	return tracing.Error(span, err)
	// }

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
			category_id int,
			foreign key(category_id) references rebrickable_part_categories(id)
		);
		
		create table if not exists rebrickable_inventories (
			id int primary key,
			version int not null,
			set_num text not null
		);

		create table if not exists rebrickable_inventory_parts (
			inventory_id int,
			part_num text not null,
			color_id int not null,
			quantity int not null,
			spare int not null,
			image_url text,
			foreign key(inventory_id) references rebrickable_inventories(id),
			foreign key(part_num) references rebrickable_parts(part_num),
			foreign key(color_id) references rebrickable_colors(id)
		);
	`

	if _, err := tx.ExecContext(ctx, createTables); err != nil {
		return tracing.Error(span, err)
	}

	clearTables := `
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
