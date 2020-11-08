package eventstore

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type EventMeta struct {
	Timestamp       time.Time `json:"meta_timestamp"`
	ID              uuid.UUID `json:"meta_id"`
	AggregateRootID uuid.UUID `json:"meta_aggregate_id"`
	Sequence        int       `json:"meta_sequence"`
	Type            string    `json:"meta_type"`
	EventVersion    int       `json:"meta_version"`
}

func (e *EventMeta) Meta() *EventMeta { return e }

type Event interface {
	Meta() *EventMeta
}
