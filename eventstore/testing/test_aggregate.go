package testing

import (
	"brickrecon/eventstore"

	uuid "github.com/satori/go.uuid"
)

type TestAggregate struct {
	*eventstore.Aggregator

	Name string
}

func BlankTestAggregate() *TestAggregate {
	a := TestAggregate{}
	a.Aggregator = eventstore.NewAggregator(a.on)
	return &a
}

func NewTestAggregate(name string) *TestAggregate {
	a := BlankTestAggregate()
	a.Apply(&TestAggregateCreated{NewID: uuid.NewV4(), Name: name})

	return a
}

func (a *TestAggregate) Rename(newName string) {
	if newName != a.Name {
		a.Apply(&TestAggregateRenamed{NewName: newName})
	}
}

func (a *TestAggregate) on(event eventstore.Event) {

	switch e := event.(type) {

	case *TestAggregateCreated:
		a.SetID(e.NewID)
		a.Name = e.Name

	case *TestAggregateRenamed:
		a.Name = e.NewName
	}

}

type TestAggregateCreated struct {
	eventstore.EventMeta

	NewID uuid.UUID
	Name  string
}

type TestAggregateRenamed struct {
	eventstore.EventMeta

	NewName string
}

type TestEvent struct {
	eventstore.EventMeta

	Name      string
	SetNumber int
}
