package eventstore

import (
	"context"

	uuid "github.com/satori/go.uuid"
)

type Backend interface {
	NewEventReader(registry map[string]Initialiser, ctx context.Context) (EventReader, error)
	NewEventWriter() EventWriter
	NewView(name string) View
}

type EventReader interface {
	Close() error
	ReadAll() bool
	ReadFor(uuid uuid.UUID) bool
	ReadFrom(offset int) bool
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
