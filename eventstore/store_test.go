package eventstore

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"path"
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

	eventOne := TestEvent{Name: "One", SetNumber: 1234}
	eventTwo := TestEvent{Name: "Two", SetNumber: 5678}

	err := es.Write(eventOne, eventTwo)
	assert.NoError(t, err)

	count, err := countLines(path.Join(es.root, "events"))
	assert.NoError(t, err)
	assert.Equal(t, 2, count)
}

func countLines(path string) (int, error) {
	r, err := os.Open(path)
	if err != nil {
		return 0, err
	}
	defer r.Close()

	buf := make([]byte, 32*1024)
	count := 0
	lineSep := []byte{'\n'}

	for {
		c, err := r.Read(buf)
		count += bytes.Count(buf[:c], lineSep)

		switch {
		case err == io.EOF:
			return count, nil

		case err != nil:
			return count, err
		}
	}
}
