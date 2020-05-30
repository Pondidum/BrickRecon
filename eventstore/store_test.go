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

	es := CreateEventStore(temp)
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

	es := CreateEventStore(temp)
	es.RegisterEvent(func() interface{} { return &TestEvent{} })

	projection := &testProjection{names: map[string]bool{}}
	es.RegisterProjection(
		"names",
		func() interface{} { return &map[string]bool{} },
		projection.Project)

	err := es.Write(TestEvent{Name: "One"})
	assert.NoError(t, err)

	assert.Contains(t, projection.names, "One")

}

func TestReadOffset(t *testing.T) {
	temp, _ := ioutil.TempDir(".", "es")
	defer func() {
		os.RemoveAll(temp)
	}()

	es := CreateEventStore(temp)
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

type testProjection struct {
	names map[string]bool
}

func (p *testProjection) Project(e interface{}) interface{} {

	switch event := e.(type) {
	case TestEvent:
		p.names[event.Name] = true
	}

	return p.names
}
