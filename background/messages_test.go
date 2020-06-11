package background

import (
	"brickrecon/eventstore"
	"brickrecon/eventstore/backend/fs"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCacheCreation(t *testing.T) {

	temp, _ := ioutil.TempDir(".", "es")
	defer func() {
		os.RemoveAll(temp)
	}()

	be, _ := fs.NewFileSystemBackend(temp)
	es := eventstore.NewEventStore(be)
	ImageCacheEvents(es.RegisterEvent)

	ic, err := loadCache(es)
	assert.NoError(t, err)
	assert.Equal(t, cacheID, ic.AggregateID())

	fromStore := blankImageCache("./app/static/img/parts")
	assert.NoError(t, es.LoadAggregate(cacheID, fromStore))
}

func TestCacheAlreadyExists(t *testing.T) {

	temp, _ := ioutil.TempDir(".", "es")
	defer func() {
		os.RemoveAll(temp)
	}()

	be, _ := fs.NewFileSystemBackend(temp)
	es := eventstore.NewEventStore(be)
	ImageCacheEvents(es.RegisterEvent)

	_, err := loadCache(es)
	assert.NoError(t, err)

	ic, err := loadCache(es)
	assert.NoError(t, err)
	assert.Equal(t, cacheID, ic.AggregateID())
}
