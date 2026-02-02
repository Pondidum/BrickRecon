package storage

import (
	"brickrecon/domain"
	"brickrecon/goes"
	"brickrecon/tracing"
	"context"
	"database/sql"

	"github.com/google/uuid"
)

type Client struct {
	db *sql.DB
	es goes.Store
}

func NewClient(ctx context.Context, dbPath string) (*Client, error) {
	ctx, span := tr.Start(ctx, "new_client")
	defer span.End()

	writer, err := Writer(ctx, dbPath)
	if err != nil {
		return nil, tracing.Error(span, err)
	}

	eventStore := goes.NewSqliteStore(writer)
	eventStore.RegisterAggregate(goes.FactoryFor(domain.BlankProject))
	eventStore.RegisterProjection(&goes.AutoProjection{})

	if err := eventStore.Initialise(ctx); err != nil {
		return nil, tracing.Error(span, err)
	}

	if err := createTables(ctx, writer); err != nil {
		return nil, tracing.Error(span, err)
	}

	client := &Client{
		db: writer,
		es: eventStore,
	}
	// change to something that also accesses the plain tables too later
	return client, nil
}

func (c *Client) BeginTx(ctx context.Context) (*sql.Tx, error) {
	return c.db.BeginTx(ctx, nil)
}

func (c *Client) LoadAggregate(ctx context.Context, aggregateID uuid.UUID, aggregate goes.Aggregate) error {
	return c.es.Load(ctx, aggregateID, aggregate)
}

func (c *Client) SaveAggregate(ctx context.Context, aggregate goes.Aggregate) error {
	return c.es.Save(ctx, aggregate)
}

func createTables(ctx context.Context, writer *sql.DB) error {
	// parts (id, name, ...)
	//
	// sets_parts (id, set_id, part_id, color, quantity)
	//
	// sets (id, name, ...)

	stmt := `
		create table if not exists parts(
			id text primary key,
			name text,
			data jsonb
		);

		create table if not exists sets(
			id text primary key,
			name text,
			data jsonb
		);
			
		create table if not exists sets_parts(
			id text primary key,
			set_id text,
			part_id text,
			color string,
			quantity int,
			foreign key(set_id) references sets(id),
			foreign key(part_id) references parts(id)
		);
	`

	_, err := writer.ExecContext(ctx, stmt)
	return err
}
