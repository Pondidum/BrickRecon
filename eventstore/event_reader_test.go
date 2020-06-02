package eventstore

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"testing"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

var testEvents string = strings.TrimSpace(`
{"Timestamp":"2020-05-31T19:35:06.615025231+03:00","ID":"7e0b10e5-2c81-4ac7-a29f-4781fd7d4e0c","AggregateRootID":"bf3faa6d-5b3f-403d-bf4f-9f7ceff972f6","Type":"TestEvent","Content":{"Name":"One","SetNumber":1234}}
{"Timestamp":"2020-05-31T19:35:06.615025231+03:00","ID":"afabeb6b-045d-489b-b085-df1a60653b1e","AggregateRootID":"43354630-8276-49d2-a30b-30710d047127","Type":"TestEvent","Content":{"Name":"One","SetNumber":1234}}
{"Timestamp":"2020-05-31T19:35:06.615025231+03:00","ID":"983f19c8-268f-4903-b4bf-37a3031d242b","AggregateRootID":"bf3faa6d-5b3f-403d-bf4f-9f7ceff972f6","Type":"TestEvent","Content":{"Name":"Two","SetNumber":1234}}
{"Timestamp":"2020-05-31T19:35:06.615025231+03:00","ID":"42b70c6a-2e38-42a8-9fa2-778ffc963c93","AggregateRootID":"bf3faa6d-5b3f-403d-bf4f-9f7ceff972f6","Type":"TestEvent","Content":{"Name":"Three","SetNumber":1234}}
{"Timestamp":"2020-05-31T19:35:06.615025231+03:00","ID":"f02227a0-ab65-4c9d-b271-a2287cb0ecf6","AggregateRootID":"81f30eee-0861-4db2-b512-3e0a6c2fdd13","Type":"TestEvent","Content":{"Name":"One","SetNumber":1234}}
`)

func TestReadingAllEvents(t *testing.T) {

	seenEvents, err := readEvents(func(reader *EventReader) bool {
		return reader.ReadAll()
	})

	assert.NoError(t, err)
	assert.Equal(t,
		[]string{
			"7e0b10e5-2c81-4ac7-a29f-4781fd7d4e0c",
			"afabeb6b-045d-489b-b085-df1a60653b1e",
			"983f19c8-268f-4903-b4bf-37a3031d242b",
			"42b70c6a-2e38-42a8-9fa2-778ffc963c93",
			"f02227a0-ab65-4c9d-b271-a2287cb0ecf6",
		},
		seenEvents,
	)
}

func TestReadingAggregateEvents(t *testing.T) {

	aggregateID := uuid.FromStringOrNil("bf3faa6d-5b3f-403d-bf4f-9f7ceff972f6")

	seenEvents, err := readEvents(func(reader *EventReader) bool {
		return reader.ReadFor(aggregateID)
	})

	assert.NoError(t, err)
	assert.Equal(t,
		[]string{
			"7e0b10e5-2c81-4ac7-a29f-4781fd7d4e0c",
			"983f19c8-268f-4903-b4bf-37a3031d242b",
			"42b70c6a-2e38-42a8-9fa2-778ffc963c93",
		},
		seenEvents,
	)
}

func TestReadingFromOffset(t *testing.T) {

	offset := 2

	seenEvents, err := readEvents(func(reader *EventReader) bool {
		return reader.ReadFrom(offset)
	})

	assert.NoError(t, err)
	assert.Equal(t,
		[]string{
			"983f19c8-268f-4903-b4bf-37a3031d242b",
			"42b70c6a-2e38-42a8-9fa2-778ffc963c93",
			"f02227a0-ab65-4c9d-b271-a2287cb0ecf6",
		},
		seenEvents,
	)
}

func TestReadingAllEventsFromOffsetZero(t *testing.T) {

	seenEvents, err := readEvents(func(reader *EventReader) bool {
		return reader.ReadFrom(0)
	})

	assert.NoError(t, err)
	assert.Equal(t,
		[]string{
			"7e0b10e5-2c81-4ac7-a29f-4781fd7d4e0c",
			"afabeb6b-045d-489b-b085-df1a60653b1e",
			"983f19c8-268f-4903-b4bf-37a3031d242b",
			"42b70c6a-2e38-42a8-9fa2-778ffc963c93",
			"f02227a0-ab65-4c9d-b271-a2287cb0ecf6",
		},
		seenEvents,
	)
}

func createTestReader(temp string) (*EventReader, error) {
	eventsFile := path.Join(temp, "events")
	ioutil.WriteFile(eventsFile, []byte(testEvents), 0666)

	reader, err := NewEventReader(
		map[string]Initialiser{
			"TestEvent": func() interface{} { return &TestEvent{} },
		},
		eventsFile,
	)

	return reader, err
}

func readEvents(method func(reader *EventReader) bool) ([]string, error) {

	temp, _ := ioutil.TempDir(".", "er")
	defer func() {
		os.RemoveAll(temp)
	}()

	reader, err := createTestReader(temp)
	if err != nil {
		return nil, err
	}

	seenEvents := []string{}

	for method(reader) {
		if event, err := reader.Event(); err != nil {
			return nil, err
		} else {
			seenEvents = append(seenEvents, event.event().ID.String())
		}
	}

	return seenEvents, nil
}

func TestEventEmbed(t *testing.T) {

	e := UserDefined{}
	e.ID = "im id"
	e.Name = "name"

	b, err := serialize(&e)

	assert.NoError(t, err)
	assert.Equal(t, `{"ID":"test","Name":"name"}`, string(b))
}

func serialize(e HasEmbed) ([]byte, error) {
	extra := e.getEmbed()
	extra.ID = "test"

	return json.Marshal(e)
}

type Embed struct {
	ID string
}

func (e *Embed) getEmbed() *Embed { return e }

type HasEmbed interface {
	getEmbed() *Embed
}

type UserDefined struct {
	Embed

	Name string
}
