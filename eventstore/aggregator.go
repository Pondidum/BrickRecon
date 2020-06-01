package eventstore

import uuid "github.com/satori/go.uuid"

type Aggregator struct {
	id      uuid.UUID
	changes []interface{}
	version int

	onEvent func(event interface{})
}

func NewAggregator(onEvent func(event interface{})) *Aggregator {
	return &Aggregator{
		onEvent: onEvent,
	}
}

func (a *Aggregator) Apply(event interface{}) {
	a.changes = append(a.changes, event)
	a.onEvent(event)
}

func (a *Aggregator) SetID(aggregateID uuid.UUID) {
	a.id = aggregateID
}

func (a *Aggregator) fromEvents(events []interface{}) {
	for _, event := range events {
		a.onEvent(event)
		a.version++
	}
}
