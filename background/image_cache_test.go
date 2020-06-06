package background

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func createState() *dto {
	return &dto{
		httpClient: &FakeClient{
			err: errors.New("nope"),
		},
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
	state.httpClient = &FakeClient{
		statusCode: 500,
		content:    []byte("some image"),
	}

	Run(state)

	assert.Len(t, state.events, 1, print(state))
	assert.IsType(t, PartAttempted{}, state.events[0], print(state))
}

func TestWhenPartFetchingFailsBecauseBodyIsUnreadable(t *testing.T) {

	state := createState()
	state.httpClient = &FakeClient{
		statusCode:  200,
		content:     []byte("some image"),
		bodyBuilder: func(content []byte) io.Reader { return &ErrorReader{} },
	}

	Run(state)

	assert.Len(t, state.events, 1, print(state))
	assert.IsType(t, PartAttempted{}, state.events[0], print(state))
}

func TestWhenPartSavingFails(t *testing.T) {

	state := createState()

	state.httpClient = &FakeClient{
		statusCode: 200,
		content:    []byte("some image"),
	}

	Run(state)

	assert.Len(t, state.events, 1, print(state))
	assert.IsType(t, PartAttempted{}, state.events[0], print(state))
}

func TestWhenPartImageDoesntExist(t *testing.T) {

	state := createState()
	state.httpClient = &FakeClient{statusCode: 404, content: []byte{}}
	state.writeFile = func(filename string, content []byte) error { return nil }

	Run(state)

	assert.Len(t, state.events, 2, print(state))
	assert.IsType(t, PartImageNotFound{}, state.events[0], print(state))
	assert.IsType(t, PartImageStored{}, state.events[1], print(state))
}

func TestWhenPartSavingWorks(t *testing.T) {

	state := createState()

	state.httpClient = &FakeClient{statusCode: 200, content: []byte("some image")}
	state.writeFile = func(filename string, content []byte) error { return nil }

	Run(state)

	assert.Len(t, state.events, 1, print(state))
	assert.IsType(t, PartImageStored{}, state.events[0], print(state))
}

func TestWhenPartFetchingExceedsMaxAttempts(t *testing.T) {

	state := createState()
	state.attempts = state.maxAttempts

	state.httpClient = &FakeClient{statusCode: 200, content: []byte("some image")}
	state.writeFile = func(filename string, content []byte) error { return nil }

	Run(state)

	assert.Len(t, state.events, 1, print(state))
	assert.IsType(t, PartFailed{}, state.events[0], print(state))
}

func print(state *dto) string {
	return strings.Join(state.transitions, " -> ")
}

type FakeClient struct {
	content     []byte
	statusCode  int
	err         error
	bodyBuilder func(content []byte) io.Reader
}

func (f *FakeClient) Do(req *http.Request) (*http.Response, error) {

	if f.bodyBuilder == nil {
		f.bodyBuilder = func(content []byte) io.Reader {
			return bytes.NewBuffer(f.content)
		}
	}

	res := &http.Response{
		StatusCode: f.statusCode,
		Body:       ioutil.NopCloser(f.bodyBuilder(f.content)),
	}

	return res, f.err
}

type ErrorReader struct{}

func (er *ErrorReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("error reading")
}
