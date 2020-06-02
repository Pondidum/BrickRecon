package eventstore

import (
	"io/ioutil"
	"os"
	"testing"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

type TestEvent struct {
	Event

	Name      string
	SetNumber int
}

func TestProjections(t *testing.T) {

	temp, _ := ioutil.TempDir(".", "es")
	defer func() {
		os.RemoveAll(temp)
	}()

	es := NewEventStore(temp)
	es.RegisterEvent(func() interface{} { return &TestEvent{} })

	es.RegisterProjection(
		"names",
		func() interface{} { return &TestProjectionState{map[string]bool{}} },
		func(state interface{}, event IsEvent) interface{} {
			m := state.(*TestProjectionState)
			e := event.(*TestEvent)

			m.Names[e.Name] = true

			return m
		})

	a := &Aggregator{
		id: uuid.NewV4(),
		changes: []IsEvent{
			&TestEvent{Name: "One"},
		},
	}
	assert.NoError(t, es.SaveAggregate(a))

	a.changes = []IsEvent{&TestEvent{Name: "Two"}}
	assert.NoError(t, es.SaveAggregate(a))

	var view TestProjectionState
	assert.NoError(t, es.ReadView("names", &view))

	assert.Contains(t, view.Names, "One")
	assert.Contains(t, view.Names, "Two")
}

type TestProjectionState struct {
	Names map[string]bool
}

func TestProjectionCatchup(t *testing.T) {
	temp, _ := ioutil.TempDir(".", "es")
	defer func() {
		os.RemoveAll(temp)
	}()

	es := NewEventStore(temp)
	es.RegisterEvent(func() interface{} { return &TestEvent{} })

	a := &Aggregator{
		id: uuid.NewV4(),
		changes: []IsEvent{
			&TestEvent{Name: "Before_1", SetNumber: 1},
			&TestEvent{Name: "Before_2", SetNumber: 2},
		},
	}

	// write some events
	assert.NoError(t, es.SaveAggregate(a))

	// register a new projection
	es.RegisterProjection(
		"logs",
		func() interface{} {
			return &OrderedEvents{}
		},
		func(state interface{}, event IsEvent) interface{} {
			m := state.(*OrderedEvents)
			e := event.(*TestEvent)

			m.Names = append(m.Names, e.Name)

			return state
		})

	// write a new event
	a.changes = []IsEvent{
		&TestEvent{Name: "After_1", SetNumber: 3},
	}
	assert.NoError(t, es.SaveAggregate(a))

	// view should contain all 3 events in order
	var view OrderedEvents
	assert.NoError(t, es.ReadView("logs", &view))

	assert.Equal(t, []string{"Before_1", "Before_2", "After_1"}, view.Names)

}

type OrderedEvents struct {
	Names []string
}

func TestAggregateSaveLoad(t *testing.T) {

	temp, _ := ioutil.TempDir(".", "er")
	defer func() {
		os.RemoveAll(temp)
	}()

	store := NewEventStore(temp)
	store.RegisterEvent(func() interface{} { return &TestAggregateCreated{} })
	store.RegisterEvent(func() interface{} { return &TestAggregateRenamed{} })

	project := NewTestAggregate("test")
	assert.NoError(t, store.SaveAggregate(project.Aggregator))

	loaded := BlankTestAggregate()
	assert.NoError(t, store.LoadAggregate(project.id, loaded.Aggregator))

	assert.Equal(t, project.id, loaded.id)
	assert.Equal(t, project.Name, loaded.Name)
	assert.Empty(t, loaded.changes)
	assert.Equal(t, 1, loaded.version)
}

func TestAggregateSave(t *testing.T) {
	temp, _ := ioutil.TempDir(".", "er")
	defer func() {
		os.RemoveAll(temp)
	}()

	store := NewEventStore(temp)
	store.RegisterEvent(func() interface{} { return &TestAggregateCreated{} })
	store.RegisterEvent(func() interface{} { return &TestAggregateRenamed{} })

	ta := NewTestAggregate("test")
	ta.Rename("two")
	assert.NoError(t, store.SaveAggregate(ta.Aggregator))
	assert.Equal(t, 2, ta.version)
	assert.Empty(t, ta.changes)

	ta.Rename("three")
	ta.Rename("four")
	assert.NoError(t, store.SaveAggregate(ta.Aggregator))
	assert.Equal(t, 4, ta.version)

}

// ------------------------------------------------------------------------- //
type TestAggregate struct {
	*Aggregator

	Name string
}

func BlankTestAggregate() *TestAggregate {
	a := TestAggregate{}
	a.Aggregator = NewAggregator(a.on)
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

func (a *TestAggregate) on(event IsEvent) {

	switch e := event.(type) {

	case *TestAggregateCreated:
		a.SetID(e.NewID)
		a.Name = e.Name

	case *TestAggregateRenamed:
		a.Name = e.NewName
	}

}

type TestAggregateCreated struct {
	Event

	NewID uuid.UUID
	Name  string
}

type TestAggregateRenamed struct {
	Event

	NewName string
}
