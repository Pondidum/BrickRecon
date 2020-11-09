package eventstore

type Aggregator struct {
	id       AggregateID
	changes  []Event
	sequence int

	onEvent func(event Event)
}

func (a *Aggregator) aggregator() *Aggregator {
	return a
}

type Aggregate interface {
	aggregator() *Aggregator
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

func (a *Aggregator) AggregateID() AggregateID {
	return a.id
}

func (a *Aggregator) SetID(aggregateID AggregateID) {
	a.id = aggregateID
}

func (a *Aggregator) fromEvents(events []Event) {
	for _, event := range events {
		a.onEvent(event)
		a.sequence++
	}
}

// these are deliberately not exposed directly on Aggregator

func ReadChanges(a Aggregate) []Event {
	return a.aggregator().changes
}
