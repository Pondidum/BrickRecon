package eventstore

import (
	"bufio"
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

func readLines(path string) ([][]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lines := [][]byte{}

	for scanner.Scan() {
		lines = append(lines, scanner.Bytes())
	}

	return lines, nil
}

func TestProjections(t *testing.T) {

	temp, _ := ioutil.TempDir(".", "es")
	defer func() {
		// os.RemoveAll(temp)
	}()

	es := CreateEventStore(temp)
	es.RegisterEvent(func() interface{} { return &TestEvent{} })

	projection := &testProjection{names: map[string]bool{}}
	es.RegisterProjection("names", projection.Project)

	err := es.Write(TestEvent{Name: "One"})
	assert.NoError(t, err)

	assert.Contains(t, projection.names, "One")

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
