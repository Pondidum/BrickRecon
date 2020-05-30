package eventstore

import (
	"io/ioutil"
	"os"
	"testing"

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

	eventOne := TestEvent{Name: "One", SetNumber: 1234}

	err := es.Write(eventOne)
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

	err := es.Write(TestEvent{Name: "One"})
	assert.NoError(t, err)

	err = es.Write(TestEvent{Name: "Two"})
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

	events := make([]interface{}, 10)
	for i := range events {
		events[i] = TestEvent{SetNumber: i}
	}

	assert.NoError(t, es.WriteEvents(events))

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

	// write some events
	assert.NoError(t, es.Write(TestEvent{Name: "Before_1", SetNumber: 1}))
	assert.NoError(t, es.Write(TestEvent{Name: "Before_2", SetNumber: 2}))

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
	assert.NoError(t, es.Write(TestEvent{Name: "After_1", SetNumber: 3}))

	// view should contain all 3 events in order
	var view OrderedEvents
	assert.NoError(t, es.ReadView("logs", &view))

	assert.Equal(t, []string{"Before_1", "Before_2", "After_1"}, view.Names)

}

type OrderedEvents struct {
	Names []string
}
