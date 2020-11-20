package eventstore

import (
	"context"

	uuid "github.com/satori/go.uuid"
)

type AggregateID string

func NewAggregateID() AggregateID {
	return AggregateID(uuid.NewV4().String())
}

type EventMiddleware func(context.Context, Event) Event

type Backend interface {
	NewEventReader(ctx context.Context, registry *EventRegistry, aggregateID AggregateID) (EventReader, error)
	NewEventWriter() EventWriter
	NewView(name string) View
	DestroyViews() error

	AllAggregates() ([]AggregateID, error)
}

type EventReader interface {
	Close() error
	Read() bool
	Event() (Event, error)
}

type EventWriter interface {
	WriteEvents(ctx context.Context, aggregateID AggregateID, changes []Event) (int, error)
}

type View interface {
	LastEventIndex(ctx context.Context) (int, error)
	ReadView(ctx context.Context, view interface{}) error
	WriteView(ctx context.Context, view interface{}, lastIndex int) error
}
