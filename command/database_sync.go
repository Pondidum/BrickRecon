package command

import (
	"brickrecon/config"
	"brickrecon/ldraw"
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
	"rebrickable_colors": map[string]string{
		"transparent": "is_trans",
	},
	"rebrickable_part_categories": map[string]string{},
	"rebrickable_parts": map[string]string{
		"category_id": "part_cat_id",
	},
	"rebrickable_inventories": map[string]string{},
	"rebrickable_sets":        map[string]string{},
	"rebrickable_inventory_parts": map[string]string{
		"spare":     "is_spare",
		"image_url": "img_url",
	},
	"rebrickable_part_relationships": map[string]string{},
	"ldraw_moved_parts":              map[string]string{},
}

func NewDatabaseSyncCommand() *DatabaseSyncCommand {
	return &DatabaseSyncCommand{
		tr: otel.Tracer("command.database.sync"),
	}
}

type DatabaseSyncCommand struct {
	tr     trace.Tracer
	dryrun bool
	debug  bool
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

	if err := writer.PrepareTables(ctx, source); err != nil {
		return tracing.Error(span, err)
	}

	switch source {
	case "ldraw":
		if err := c.syncLDraw(ctx, writer); err != nil {
			return tracing.Error(span, err)
		}

	case "rebrickable":

		for tableName := range ColumnSources {
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

	res, err := http.Get("https://library.ldraw.org/library/updates/complete.zip")
	if err != nil {
		return tracing.Error(span, err)
	}
	defer res.Body.Close()

	parts, err := ldraw.ParseDatabaseArchive(ctx, res.Body)
	if err != nil {
		return tracing.Error(span, err)
	}

	insert, done, err := writer.Prepare(ctx, "ldraw_moved_parts", []string{"old_part_num", "new_part_num"})
	if err != nil {
		return tracing.Error(span, err)
	}
	defer done()

	for old, new := range parts {
		if new != "" {

			if err := insert([]string{old, new}); err != nil {
				return tracing.Error(span, err)
			}
		}
	}

	return nil
}

type writer interface {
	PrepareTables(ctx context.Context, source string) error
	Cancel() error
	Prepare(ctx context.Context, tableName string, header []string) (func(record []string) error, func() error, error)
	Finish(ctx context.Context) error
}

type dbwriter struct {
	client *storage.Client
	tx     *sql.Tx
}

func NewDbWriter(ctx context.Context, db *storage.Client) (writer, error) {

	tx, err := db.BeginTx(ctx)
	if err != nil {
		return nil, err
	}

	return &dbwriter{client: db, tx: tx}, nil
}

func (db *dbwriter) PrepareTables(ctx context.Context, source string) error {
	switch source {
	case "ldraw":
		return db.PrepareLDrawTables(ctx)
	case "rebrickable":
		return db.PrepareRebrickableTables(ctx)
	default:
		return fmt.Errorf("unsupported source %s", source)
	}
}

func (db *dbwriter) PrepareRebrickableTables(ctx context.Context) error {

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

		create table if not exists rebrickable_part_relationships (
			rel_type text not null,
			child_part_num text not null references rebrickable_parts(part_num) deferrable initially deferred,
			parent_part_num text not null references rebrickable_parts(part_num) deferrable initially deferred
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

	fmt.Println("Creating missing tables...")
	if _, err := db.tx.ExecContext(ctx, createTables); err != nil {
		return err
	}

	clearTables := `
		delete from rebrickable_sets where 1=1;
		delete from rebrickable_inventory_parts where 1=1;
		delete from rebrickable_inventories where 1=1;
		delete from rebrickable_parts where 1=1;
		delete from rebrickable_part_relationships where 1=1;
		delete from rebrickable_colors where 1=1;
		delete from rebrickable_part_categories where 1=1;
	`

	fmt.Println("Clearing old data...")
	if _, err := db.tx.ExecContext(ctx, clearTables); err != nil {
		return err
	}

	fmt.Println("Tables created")

	return nil
}

func (db *dbwriter) PrepareLDrawTables(ctx context.Context) error {

	createTables := `
		create table if not exists ldraw_moved_parts (
			old_part_num text,
			new_part_num text
		);
	`

	fmt.Println("Creating missing tables...")
	if _, err := db.tx.ExecContext(ctx, createTables); err != nil {
		return err
	}

	clearTables := `
		delete from ldraw_moved_parts where 1=1;
	`

	fmt.Println("Clearing old data...")
	if _, err := db.tx.ExecContext(ctx, clearTables); err != nil {
		return err
	}

	fmt.Println("Tables created")

	return nil
}

func (db *dbwriter) Cancel() error {
	return db.tx.Rollback()
}

func (db *dbwriter) Prepare(ctx context.Context, tableName string, header []string) (func(record []string) error, func() error, error) {

	columnMapping, err := db.createColumnMapping(ctx, tableName)
	if err != nil {
		return nil, nil, err
	}

	stmt, err := db.tx.PrepareContext(ctx, db.generateInsertSql(tableName, header, columnMapping))
	if err != nil {
		return nil, nil, err
	}

	count := 0

	insert := func(record []string) error {

		args := make([]any, 0, len(columnMapping))
		for i, name := range header {
			if _, found := columnMapping[name]; found {
				args = append(args, record[i])
			}
		}

		if _, err := stmt.ExecContext(ctx, args...); err != nil {
			return err
		}

		count++

		return nil
	}

	done := func() error {
		fmt.Println("Inserted", count, tableName, "records")
		return stmt.Close()
	}

	return insert, done, nil
}

func (db *dbwriter) generateInsertSql(tableName string, csvHeaders []string, columnMapping map[string]string) string {

	sb := strings.Builder{}
	sb.WriteString("insert into ")
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

func (db *dbwriter) createColumnMapping(ctx context.Context, tableName string) (map[string]string, error) {

	row, err := db.tx.QueryContext(ctx,
		"select name from pragma_table_info(@table)",
		sql.Named("table", tableName))
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

func (db *dbwriter) Finish(ctx context.Context) error {
	fmt.Println("Committing transaction...")
	if err := db.tx.Commit(); err != nil {
		return err
	}

	if err := db.client.Vacuum(ctx); err != nil {
		return err
	}

	if err := db.client.Close(ctx); err != nil {
		return err
	}

	return nil
}

var _ writer = &drywriter{}

type drywriter struct {
	debug    bool
	finished bool
}

func (dw *drywriter) PrepareTables(ctx context.Context, source string) error {
	fmt.Println("Dry run starting...")
	return nil
}

func (dw *drywriter) Cancel() error {
	if !dw.finished {
		fmt.Println("Dry run cancelled")
	}
	return nil
}

func (dw *drywriter) Prepare(ctx context.Context, tableName string, header []string) (func(record []string) error, func() error, error) {
	counter := 0
	if dw.debug {
		fmt.Println("header", header)
	}
	return func(record []string) error {
			if dw.debug && counter < 5 {
				fmt.Println("insert", record)
			}
			counter++
			return nil
		}, func() error {
			fmt.Println("Would have inserted", counter, tableName, "records")
			return nil
		}, nil
}

func (dw *drywriter) Finish(ctx context.Context) error {
	dw.finished = true
	fmt.Println("Dry run finished")
	return nil
}
