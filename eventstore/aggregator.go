package eventstore

import uuid "github.com/satori/go.uuid"

type Aggregator struct {
	id      uuid.UUID
	changes []IsEvent
	version int

	onEvent func(event IsEvent)
}

func NewAggregator(onEvent func(event IsEvent)) *Aggregator {
	return &Aggregator{
		onEvent: onEvent,
	}
}

func (a *Aggregator) Apply(event IsEvent) {
	a.changes = append(a.changes, event)
	a.onEvent(event)
}

func (a *Aggregator) SetID(aggregateID uuid.UUID) {
	a.id = aggregateID
}

func (a *Aggregator) fromEvents(events []IsEvent) {
	for _, event := range events {
		a.onEvent(event)
		a.version++
	}
}
