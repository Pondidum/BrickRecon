package background

import (
	"errors"
	"io"
	"strings"
	"testing"

	"mvc/testutil"

	"github.com/stretchr/testify/assert"
)

func createState() *dto {
	return &dto{
		httpClient: testutil.HttpNetworkErrorClient(),
		writeFile: func(filename string, content []byte) error {
			return errors.New("File write failed")
		},

		invalidImage: []byte("invalid image placeholder"),

		attempts:    0,
		maxAttempts: 5,
		format:      "png",

		partID:   "567b",
		colourID: 85,
	}

}

func TestWhenPartFetchingFails(t *testing.T) {

	state := createState()

	Run(state)

	assert.Len(t, state.events, 1, print(state))
	assert.IsType(t, PartAttempted{}, state.events[0], print(state))
}

func TestWhenPartFetchingFailsBecauseServerErrors(t *testing.T) {

	state := createState()
	state.httpClient = testutil.HttpServerErrorClient()

	Run(state)

	assert.Len(t, state.events, 1, print(state))
	assert.IsType(t, PartAttempted{}, state.events[0], print(state))
}

func TestWhenPartFetchingFailsBecauseBodyIsUnreadable(t *testing.T) {

	state := createState()
	state.httpClient = &testutil.FakeClient{
		StatusCode:  200,
		BodyBuilder: func(content []byte) io.Reader { return testutil.NewErrorReader() },
	}

	Run(state)

	assert.Len(t, state.events, 1, print(state))
	assert.IsType(t, PartAttempted{}, state.events[0], print(state))
}

func TestWhenPartSavingFails(t *testing.T) {

	state := createState()
	state.httpClient = testutil.HttpOkClient([]byte("some image"))

	Run(state)

	assert.Len(t, state.events, 1, print(state))
	assert.IsType(t, PartAttempted{}, state.events[0], print(state))
}

func TestWhenPartImageDoesntExist(t *testing.T) {

	state := createState()
	state.httpClient = testutil.HttpNotFoundClient()
	state.writeFile = func(filename string, content []byte) error { return nil }

	Run(state)

	assert.Len(t, state.events, 2, print(state))
	assert.IsType(t, PartImageNotFound{}, state.events[0], print(state))
	assert.IsType(t, PartImageStored{}, state.events[1], print(state))
}

func TestWhenPartSavingWorks(t *testing.T) {

	state := createState()

	state.httpClient = testutil.HttpOkClient([]byte("some image"))
	state.writeFile = func(filename string, content []byte) error { return nil }

	Run(state)

	assert.Len(t, state.events, 1, print(state))
	assert.IsType(t, PartImageStored{}, state.events[0], print(state))
}

func TestWhenPartFetchingExceedsMaxAttempts(t *testing.T) {

	state := createState()
	state.attempts = state.maxAttempts

	Run(state)

	assert.Len(t, state.events, 1, print(state))
	assert.IsType(t, PartFailed{}, state.events[0], print(state))
}

func print(state *dto) string {
	return strings.Join(state.transitions, " -> ")
}
