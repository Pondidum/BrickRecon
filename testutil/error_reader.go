package testutil

import "errors"

type ErrorReader struct {
	err error
}

func (r *ErrorReader) Read(p []byte) (n int, err error) {
	return 0, r.err
}

func NewErrorReader() *ErrorReader {
	return &ErrorReader{
		err: errors.New("Error reading stream"),
	}
}
