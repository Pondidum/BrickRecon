package fs

import (
	"brickrecon/eventstore"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"testing"
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

var testEvents string = strings.TrimSpace(`
{"meta_timestamp":"2020-05-31T19:35:06.615025231+03:00","meta_id":"7e0b10e5-2c81-4ac7-a29f-4781fd7d4e0c","meta_aggregate_id":"bf3faa6d-5b3f-403d-bf4f-9f7ceff972f6","meta_version":4, "meta_type":"TestEvent","Content":{"Name":"One","SetNumber":1234}}
{"meta_timestamp":"2020-05-31T19:35:06.615025231+03:00","meta_id":"afabeb6b-045d-489b-b085-df1a60653b1e","meta_aggregate_id":"43354630-8276-49d2-a30b-30710d047127","meta_version":0, "meta_type":"TestEvent","Content":{"Name":"One","SetNumber":1234}}
{"meta_timestamp":"2020-05-31T19:35:06.615025231+03:00","meta_id":"983f19c8-268f-4903-b4bf-37a3031d242b","meta_aggregate_id":"bf3faa6d-5b3f-403d-bf4f-9f7ceff972f6","meta_version":5, "meta_type":"TestEvent","Content":{"Name":"Two","SetNumber":1234}}
{"meta_timestamp":"2020-05-31T19:35:06.615025231+03:00","meta_id":"42b70c6a-2e38-42a8-9fa2-778ffc963c93","meta_aggregate_id":"bf3faa6d-5b3f-403d-bf4f-9f7ceff972f6","meta_version":6, "meta_type":"TestEvent","Content":{"Name":"Three","SetNumber":1234}}
{"meta_timestamp":"2020-05-31T19:35:06.615025231+03:00","meta_id":"f02227a0-ab65-4c9d-b271-a2287cb0ecf6","meta_aggregate_id":"81f30eee-0861-4db2-b512-3e0a6c2fdd13","meta_version":0, "meta_type":"TestEvent","Content":{"Name":"One","SetNumber":1234}}
`)

func TestDeserialization(t *testing.T) {

	temp, _ := ioutil.TempDir(".", "er")
	defer func() {
		os.RemoveAll(temp)
	}()

	reader, err := createTestReader(temp)
	assert.NoError(t, err)

	reader.ReadAll()
	event, err := reader.Event()
	assert.NoError(t, err)

	expectedAggregateID := uuid.Must(uuid.FromString("bf3faa6d-5b3f-403d-bf4f-9f7ceff972f6"))
	expectedEventID := uuid.Must(uuid.FromString("7e0b10e5-2c81-4ac7-a29f-4781fd7d4e0c"))
	expectedTime, err := time.Parse(time.RFC3339Nano, "2020-05-31T19:35:06.615025231+03:00")

	assert.NoError(t, err)

	assert.Equal(t, expectedAggregateID, event.AggregateID())
	assert.Equal(t, expectedEventID, event.Meta().ID)
	assert.Equal(t, expectedAggregateID, event.Meta().AggregateRootID)
	assert.Equal(t, "TestEvent", event.Meta().Type)
	assert.Equal(t, 4, event.Meta().Version)
	assert.Equal(t, expectedTime, event.Meta().Timestamp)
}

func TestReadingAllEvents(t *testing.T) {

	seenEvents, err := readEvents(func(reader *FsEventReader) bool {
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

	seenEvents, err := readEvents(func(reader *FsEventReader) bool {
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

	seenEvents, err := readEvents(func(reader *FsEventReader) bool {
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

	seenEvents, err := readEvents(func(reader *FsEventReader) bool {
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

func createTestReader(temp string) (*FsEventReader, error) {
	eventsFile := path.Join(temp, "events")
	ioutil.WriteFile(eventsFile, []byte(testEvents), 0666)

	reader, err := NewEventReader(
		map[string]eventstore.Initialiser{
			"TestEvent": func() interface{} { return &TestEvent{} },
		},
		eventsFile,
	)

	return reader, err
}

func readEvents(method func(reader *FsEventReader) bool) ([]string, error) {

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
			seenEvents = append(seenEvents, event.Meta().ID.String())
		}
	}

	return seenEvents, nil
}

type TestEvent struct {
	eventstore.EventMeta

	Name      string
	SetNumber int
}
