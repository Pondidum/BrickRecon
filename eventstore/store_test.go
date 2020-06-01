package eventstore

import (
	"io/ioutil"
	"mvc/lego"
	"os"
	"testing"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

type TestEvent struct {
	Name      string
	SetNumber int
}

func TestWritingEvents(t *testing.T) {

	temp, _ := ioutil.TempDir(".", "es")
	defer func() {
		os.RemoveAll(temp)
	}()

	es := NewEventStore(temp)
	es.RegisterEvent(func() interface{} { return &TestEvent{} })

	id := uuid.NewV4()
	eventOne := TestEvent{Name: "One", SetNumber: 1234}

	err := es.Write(id, eventOne)
	assert.NoError(t, err)

	events, err := es.ReadEvents(0)
	assert.NoError(t, err)
	assert.Len(t, events, 1)

	for _, e := range events {

		switch event := e.(type) {
		case *TestEvent:
			assert.Equal(t, "One", event.Name)
		default:
			assert.Fail(t, "")
		}

	}
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
		func(state, event interface{}) interface{} {
			m := state.(*TestProjectionState)
			e := event.(*TestEvent)

			m.Names[e.Name] = true

			return m
		})

	id := uuid.NewV4()
	err := es.Write(id, TestEvent{Name: "One"})
	assert.NoError(t, err)

	err = es.Write(id, TestEvent{Name: "Two"})
	assert.NoError(t, err)

	var view TestProjectionState
	err = es.ReadView("names", &view)
	assert.NoError(t, err)

	assert.Contains(t, view.Names, "One")
	assert.Contains(t, view.Names, "Two")
}

type TestProjectionState struct {
	Names map[string]bool
}

func TestReadOffset(t *testing.T) {
	temp, _ := ioutil.TempDir(".", "es")
	defer func() {
		os.RemoveAll(temp)
	}()

	es := NewEventStore(temp)
	es.RegisterEvent(func() interface{} { return &TestEvent{} })

	id := uuid.NewV4()
	events := make([]interface{}, 10)
	for i := range events {
		events[i] = TestEvent{SetNumber: i}
	}

	assert.NoError(t, es.WriteEvents(id, events))

	readEvents, err := es.ReadEvents(7)
	assert.NoError(t, err)

	assert.Equal(t, 7, readEvents[0].(*TestEvent).SetNumber)
	assert.Len(t, readEvents, 3)
}

func TestProjectionCatchup(t *testing.T) {
	temp, _ := ioutil.TempDir(".", "es")
	defer func() {
		os.RemoveAll(temp)
	}()

	es := NewEventStore(temp)
	es.RegisterEvent(func() interface{} { return &TestEvent{} })

	id := uuid.NewV4()

	// write some events
	assert.NoError(t, es.Write(id, TestEvent{Name: "Before_1", SetNumber: 1}))
	assert.NoError(t, es.Write(id, TestEvent{Name: "Before_2", SetNumber: 2}))

	// register a new projection
	es.RegisterProjection(
		"logs",
		func() interface{} {
			return &OrderedEvents{}
		},
		func(state, event interface{}) interface{} {
			m := state.(*OrderedEvents)
			e := event.(*TestEvent)

			m.Names = append(m.Names, e.Name)

			return state
		})

	// write a new event
	assert.NoError(t, es.Write(id, TestEvent{Name: "After_1", SetNumber: 3}))

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
	store.RegisterEvent(func() interface{} { return &lego.ProjectCreated{} })
	store.RegisterEvent(func() interface{} { return &lego.PartsAdded{} })

	project := lego.NewProject("test", []lego.Part{})
	assert.NoError(t, store.SaveAggregate(project))

	var loaded lego.Project
	assert.NoError(t, store.LoadAggregate(project.ID(), &loaded))

	assert.Equal(t, project.ID(), loaded.ID())
	assert.Equal(t, project.Name, loaded.Name)
	assert.Empty(t, loaded.Changes())
	assert.Equal(t, 2, loaded.Version())
}
