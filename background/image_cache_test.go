package background

import (
	"errors"
	"io"
	"strings"
	"testing"

	"mvc/testutil"

	"github.com/stretchr/testify/assert"
)

func createState() *fsm {
	return &fsm{
		httpClient: testutil.HttpNetworkErrorClient(),
		writeFile: func(filename string, content []byte) error {
			return errors.New("File write failed")
		},

		invalidImage: []byte("invalid image placeholder"),

		attempts:    0,
		maxAttempts: 5,

		partID:   "567b",
		colourID: 85,
	}

}

func TestWhenPartFetchingFails(t *testing.T) {

	state := createState()
	state.Run()

	assert.Len(t, state.events, 1, print(state))

	event := state.events[0].(*PartAttempted)

	assert.Equal(t, state.partID, event.PartID)
	assert.Equal(t, state.colourID, event.ColourID)
	assert.Equal(t, "No network path available", event.Error)
}

func TestWhenPartFetchingFailsBecauseServerErrors(t *testing.T) {

	state := createState()
	state.httpClient = testutil.HttpServerErrorClient()
	state.Run()

	assert.Len(t, state.events, 1, print(state))
	event := state.events[0].(*PartAttempted)

	assert.Equal(t, state.partID, event.PartID)
	assert.Equal(t, state.colourID, event.ColourID)
	assert.Equal(t, "Unexpected statusCode: 500", event.Error)
}

func TestWhenPartFetchingFailsBecauseBodyIsUnreadable(t *testing.T) {

	state := createState()
	state.httpClient = &testutil.FakeClient{
		StatusCode:  200,
		BodyBuilder: func(content []byte) io.Reader { return testutil.NewErrorReader() },
	}
	state.Run()

	assert.Len(t, state.events, 1, print(state))
	event := state.events[0].(*PartAttempted)

	assert.Equal(t, state.partID, event.PartID)
	assert.Equal(t, state.colourID, event.ColourID)
	assert.Equal(t, "Error reading stream", event.Error)
}

func TestWhenPartSavingFails(t *testing.T) {

	state := createState()
	state.httpClient = testutil.HttpOkClient([]byte("some image"))
	state.Run()

	assert.Len(t, state.events, 1, print(state))
	event := state.events[0].(*PartAttempted)

	assert.Equal(t, state.partID, event.PartID)
	assert.Equal(t, state.colourID, event.ColourID)
	assert.Equal(t, "File write failed", event.Error)
}

func TestWhenPartImageDoesntExist(t *testing.T) {

	var fileWritten string
	var contentWritten []byte

	state := createState()

	state.httpClient = testutil.HttpNotFoundClient()
	state.writeFile = func(filename string, content []byte) error {
		fileWritten = filename
		contentWritten = content
		return nil
	}
	state.Run()

	attempt := state.events[0].(*PartImageNotFound)

	assert.Len(t, state.events, 2, print(state))
	assert.Equal(t, state.partID, attempt.PartID)
	assert.Equal(t, state.colourID, attempt.ColourID)

	assert.Equal(t, "567b-85.png", fileWritten)
	assert.Equal(t, "invalid image placeholder", string(contentWritten))
}

func TestWhenPartSavingWorks(t *testing.T) {

	var fileWritten string
	var contentWritten []byte

	state := createState()

	state.httpClient = testutil.HttpOkClient([]byte("some image data"))
	state.writeFile = func(filename string, content []byte) error {
		fileWritten = filename
		contentWritten = content
		return nil
	}
	state.Run()

	event := state.events[0].(*PartImageStored)

	assert.Len(t, state.events, 1, print(state))
	assert.Equal(t, state.partID, event.PartID)
	assert.Equal(t, state.colourID, event.ColourID)

	assert.Equal(t, "567b-85.png", fileWritten)
	assert.Equal(t, "some image data", string(contentWritten))

}

func TestWhenPartFetchingExceedsMaxAttempts(t *testing.T) {

	var fileWritten string
	var contentWritten []byte

	state := createState()
	state.attempts = state.maxAttempts
	state.writeFile = func(filename string, content []byte) error {
		fileWritten = filename
		contentWritten = content
		return nil
	}

	state.Run()

	exceed := state.events[0].(*PartFetchAttemptsExceeded)

	assert.Len(t, state.events, 2, print(state))
	assert.Equal(t, state.partID, exceed.PartID)
	assert.Equal(t, state.colourID, exceed.ColourID)

	assert.Equal(t, "567b-85.png", fileWritten)
	assert.Equal(t, "invalid image placeholder", string(contentWritten))
}

func print(state *fsm) string {
	return strings.Join(state.transitions, " -> ")
}
