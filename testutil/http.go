package testutil

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
)

type FakeClient struct {
	Content     []byte
	StatusCode  int
	Err         error
	BodyBuilder func(content []byte) io.Reader
}

func (f *FakeClient) Do(req *http.Request) (*http.Response, error) {

	if f.BodyBuilder == nil {
		f.BodyBuilder = func(content []byte) io.Reader {
			return bytes.NewBuffer(content)
		}
	}

	res := &http.Response{
		StatusCode: f.StatusCode,
		Body:       ioutil.NopCloser(f.BodyBuilder(f.Content)),
	}

	return res, f.Err
}

func HttpNetworkErrorClient() *FakeClient {
	return &FakeClient{
		Err: errors.New("No network path available"),
	}
}

func HttpOkClient(content []byte) *FakeClient {
	return &FakeClient{
		StatusCode: 200,
		Content:    content,
	}
}

func HttpNotFoundClient() *FakeClient {
	return &FakeClient{
		StatusCode: 404,
		Content:    []byte("404 not found"),
	}
}

func HttpServerErrorClient() *FakeClient {
	return &FakeClient{
		StatusCode: 500,
		Content:    []byte("500 internal server error"),
	}
}
