package command

import (
	"brickrecon/storage"
	"context"
	"database/sql"
	"fmt"
	"slices"
	"strings"
)

type Table struct {
	create        string
	insertStmt    string
	sourceColumns []string
}

var Tables = map[string]Table{
	"ldraw_moved_parts": {
		create:        `create table if not exists ldraw_moved_parts (old_part_num text primary_key, new_part_num text not null)`,
		sourceColumns: []string{"old_part_num", "new_part_num"},
	},
	"ldraw_alternate_ids": {
		create: `create table if not exists ldraw_alternate_ids (
			part_num text not null,
			alt_part_num text not null,
			unique (part_num, alt_part_num)
		)`,
		sourceColumns: []string{"part_num", "alt_part_num"},
		insertStmt:    `insert into ldraw_alternate_ids values (?, ?) on conflict do nothing`,
	},

	"rebrickable_part_categories ": {
		create: `create table if not exists rebrickable_part_categories (
			id int primary key,
			name text not null
		)`,
		sourceColumns: []string{"id", "name"},
	},

	"rebrickable_colors ": {
		create: `create table if not exists rebrickable_colors (
			id int primary key,
			name text not null,
			rgb text,
			transparent int
		)`,
		sourceColumns: []string{"id", "name", "rgb", "is_trans"},
	},

	"rebrickable_parts ": {
		create: `create table if not exists rebrickable_parts (
			part_num text primary key,
			name text not null,
			category_id int references rebrickable_part_categories(id) deferrable initially deferred
		)`,
		sourceColumns: []string{"part_num", "name", "part_cat_id"},
	},

	"rebrickable_part_relationships ": {
		create: `create table if not exists rebrickable_part_relationships (
			rel_type text not null,
			child_part_num text not null references rebrickable_parts(part_num) deferrable initially deferred,
			parent_part_num text not null references rebrickable_parts(part_num) deferrable initially deferred
		)`,
		sourceColumns: []string{"rel_type", "child_part_num", "parent_part_num"},
	},

	"rebrickable_sets ": {
		create: `create table if not exists rebrickable_sets (
			set_num text primary key,
			name text not null,
			year int not null,
			theme_id int not null,
			num_parts int not null
		)`,
		sourceColumns: []string{"set_num", "name", "year", "theme_id", "num_parts"},
	},

	"rebrickable_minifigs": {
		create: `create table if not exists rebrickable_minifigs(
			fig_num text primary key,
			name text not null,
			num_parts int not null
		)`,
		sourceColumns: []string{"fig_num", "name", "num_parts"},
	},

	"rebrickable_inventory_minifigs ": {
		create: `create table if not exists rebrickable_inventory_minifigs (
			inventory_id text references rebrickable_inventories(id) deferrable initially deferred,
			fig_num text references rebrickable_minifigs(fig_num) deferrable initially deferred,
			quantity int not null
		)`,
		sourceColumns: []string{"inventory_id", "fig_num", "quentity"},
	},

	"rebrickable_inventories ": {
		create: `create table if not exists rebrickable_inventories (
			id int primary key,
			version int not null,
			set_num text not null -- note, no fk here as this could also be from the minifigs table
		)`,
		sourceColumns: []string{"id", "version", "set_num"},
	},

	"rebrickable_inventory_parts ": {
		create: `create table if not exists rebrickable_inventory_parts (
			inventory_id int references rebrickable_inventories(id) deferrable initially deferred,
			part_num text not null references rebrickable_parts(part_num) deferrable initially deferred,
			color_id int not null references rebrickable_colors(id) deferrable initially deferred,
			quantity int not null,
			spare int not null,
			image_url text
		)`,
		sourceColumns: []string{"inventory_id", "part_num", "color_id", "quantity", "is_spare", "img_url"},
	},
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

func (db *dbwriter) CreateTables(ctx context.Context, source string) error {

	fmt.Print("Creating missing tables...")
	for _, table := range Tables {
		if _, err := db.tx.ExecContext(ctx, table.create); err != nil {
			return err
		}
	}
	fmt.Println("done")
	return nil
}

func (db *dbwriter) ClearTables(ctx context.Context, source string) error {
	fmt.Print("Clearing old data...")

	for name := range Tables {
		if strings.HasPrefix(name, source) {
			if _, err := db.tx.ExecContext(ctx, fmt.Sprintf(`delete from %s where 1=1`, name)); err != nil {
				return err
			}
		}
	}

	fmt.Println("done")
	return nil
}

func (db *dbwriter) Cancel() error {
	return db.tx.Rollback()
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

func (db *dbwriter) Prepare(ctx context.Context, tableName string, header []string) (func(record []string) error, func() error, error) {

	table, found := Tables[tableName]
	if !found {
		return nil, nil, fmt.Errorf("unsupported table %s", tableName)
	}

	insertStmt := table.insertStmt
	if insertStmt == "" {
		insertStmt = fmt.Sprintf(`insert into %s values (?%s)`, tableName, strings.Repeat(", ?", len(table.sourceColumns)-1))
	}

	stmt, err := db.tx.PrepareContext(ctx, insertStmt)
	if err != nil {
		return nil, nil, err
	}

	indexes := []int{}
	for _, col := range table.sourceColumns {
		if idx := slices.Index(header, col); idx != -1 {
			indexes = append(indexes, idx)
		}
	}

	insert := func(record []string) error {

		args := make([]any, 0, len(record))
		for _, idx := range indexes {
			args = append(args, record[idx])
		}

		_, err := stmt.ExecContext(ctx, args...)
		return err
	}

	done := stmt.Close

	return insert, done, nil

}
