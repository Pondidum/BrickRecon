package eventstore

import (
	"context"

	uuid "github.com/satori/go.uuid"
)

type EventMiddleware func(context.Context, Event) Event

type Backend interface {
	NewEventReader(registry *EventRegistry, aggregateID uuid.UUID) (EventReader, error)
	NewEventWriter() EventWriter
	NewView(name string) View
	DestroyViews() error

	AllAggregates() ([]uuid.UUID, error)
}

type EventReader interface {
	Close() error
	Read() bool
	Event() (Event, error)
}

type EventWriter interface {
	WriteEvents(ctx context.Context, aggregateID uuid.UUID, changes []Event) (int, error)
}

type View interface {
	LastEventIndex(ctx context.Context) (int, error)
	ReadView(ctx context.Context, view interface{}) error
	WriteView(ctx context.Context, view interface{}, lastIndex int) error
}
