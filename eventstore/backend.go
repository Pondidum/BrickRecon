package eventstore

import (
	uuid "github.com/satori/go.uuid"
)

type Backend interface {
	NewEventReader(map[string]Initialiser) (EventReader, error)
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
	WriteEvents(aggregateID uuid.UUID, currentVersion int, changes []Event) (int, error)
}

type View interface {
	LastEventIndex() (int, error)
	ReadView(view interface{}) error
	WriteView(view interface{}, lastIndex int) error
}
