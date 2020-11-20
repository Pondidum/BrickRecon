package testing

import (
	"brickrecon/eventstore"
	"brickrecon/eventstore/backend/fs"
	"context"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func createBackend() (eventstore.Backend, func()) {

	temp, _ := ioutil.TempDir(".", "es")
	be, _ := fs.NewAggregateBackend(temp)

	return be, func() { os.RemoveAll(temp) }

}

func TestEventRegistration(t *testing.T) {

	es := eventstore.NewEventStore(nil)

	assert.NoError(t, es.RegisterEvent(context.Background(), func() interface{} {
		return &TestEvent{}
	}))

	assert.Error(t, es.RegisterEvent(context.Background(), func() interface{} {
		return TestEvent{}
	}))
}

func TestProjections(t *testing.T) {

	be, cleanup := createBackend()
	defer cleanup()

	es := eventstore.NewEventStore(be)
	es.RegisterEvent(context.Background(), func() interface{} { return &TestEvent{} })

	es.RegisterProjection(
		context.Background(),
		&testProjection{
			name: "names",
			init: func() interface{} { return &TestProjectionState{map[string]bool{}} },
			project: func(ctx context.Context, state interface{}, event eventstore.Event) interface{} {
				m := state.(*TestProjectionState)
				e := event.(*TestEvent)

				m.Names[e.Name] = true

				return m
			},
		})

	a := eventstore.NewAggregator(func(e eventstore.Event) {})
	a.SetID(eventstore.NewAggregateID())
	a.Apply(&TestEvent{Name: "One"})

	assert.NoError(t, es.SaveAggregate(context.Background(), a))

	a.Apply(&TestEvent{Name: "Two"})
	assert.NoError(t, es.SaveAggregate(context.Background(), a))

	var view TestProjectionState
	assert.NoError(t, es.ReadView(context.Background(), "names", &view))

	assert.Contains(t, view.Names, "One")
	assert.Contains(t, view.Names, "Two")
}

func TestAggregateSaveLoad(t *testing.T) {
	be, cleanup := createBackend()
	defer cleanup()

	store := eventstore.NewEventStore(be)
	store.RegisterEvent(context.Background(), func() interface{} { return &TestAggregateCreated{} })
	store.RegisterEvent(context.Background(), func() interface{} { return &TestAggregateRenamed{} })

	project := NewTestAggregate("test")
	assert.NoError(t, store.SaveAggregate(context.Background(), project.Aggregator))

	loaded := BlankTestAggregate()
	assert.NoError(t, store.LoadAggregate(context.Background(), project.AggregateID(), loaded.Aggregator))

	assert.Equal(t, project.AggregateID(), loaded.AggregateID())
	assert.Equal(t, project.Name, loaded.Name)
	assert.Empty(t, eventstore.ReadChanges(loaded))
	// assert.Equal(t, 1, loaded.version)
}

func TestAggregateSave(t *testing.T) {
	be, cleanup := createBackend()
	defer cleanup()

	store := eventstore.NewEventStore(be)
	store.RegisterEvent(context.Background(), func() interface{} { return &TestAggregateCreated{} })
	store.RegisterEvent(context.Background(), func() interface{} { return &TestAggregateRenamed{} })

	ta := NewTestAggregate("test")
	ta.Rename("two")
	assert.NoError(t, store.SaveAggregate(context.Background(), ta.Aggregator))
	// assert.Equal(t, 2, ta.version)
	assert.Empty(t, eventstore.ReadChanges(ta))

	ta.Rename("three")
	ta.Rename("four")
	assert.NoError(t, store.SaveAggregate(context.Background(), ta.Aggregator))
	// assert.Equal(t, 4, ta.version)

}

func TestWhenAggregateIsntFound(t *testing.T) {
	be, cleanup := createBackend()
	defer cleanup()

	store := eventstore.NewEventStore(be)
	store.SaveAggregate(context.Background(), NewTestAggregate("test"))

	// not the same ID
	id := eventstore.NewAggregateID()

	a := BlankTestAggregate()
	err := store.LoadAggregate(context.Background(), id, a)

	assert.True(t, strings.HasPrefix(err.Error(), "No aggregate found for ID"))
}

func TestWhenReadingFromEmptyStore(t *testing.T) {
	be, cleanup := createBackend()
	defer cleanup()

	store := eventstore.NewEventStore(be)
	id := eventstore.NewAggregateID()
	a := BlankTestAggregate()
	err := store.LoadAggregate(context.Background(), id, a)

	assert.True(t, strings.HasPrefix(err.Error(), "No aggregate found for ID"))
	assert.True(t, eventstore.IsAggregateNotFound(err))
}
