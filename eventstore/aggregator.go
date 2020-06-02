package eventstore

import uuid "github.com/satori/go.uuid"

type Aggregator struct {
	id      uuid.UUID
	changes []Event
	version int

	onEvent func(event Event)
}

func NewAggregator(onEvent func(event Event)) *Aggregator {
	return &Aggregator{
		onEvent: onEvent,
	}
}

func (a *Aggregator) Apply(event Event) {
	a.changes = append(a.changes, event)
	a.onEvent(event)
}

func (a *Aggregator) SetID(aggregateID uuid.UUID) {
	a.id = aggregateID
}

func (a *Aggregator) fromEvents(events []Event) {
	for _, event := range events {
		a.onEvent(event)
		a.version++
	}
}
