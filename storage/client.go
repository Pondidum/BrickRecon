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

func (c *Client) Vacuum(ctx context.Context) error {
	_, err := c.db.ExecContext(ctx, "vacuum")
	return err
}

func (c *Client) LoadAggregate(ctx context.Context, aggregateID uuid.UUID, aggregate goes.Aggregate) error {
	return c.es.Load(ctx, aggregateID, aggregate)
}

func (c *Client) SaveAggregate(ctx context.Context, aggregate goes.Aggregate) error {
	return c.es.Save(ctx, aggregate)
}

func (c *Client) Close(ctx context.Context) error {
	return c.db.Close()
}
