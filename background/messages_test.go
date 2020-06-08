package background

import (
	"io/ioutil"
	"mvc/eventstore"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCacheCreation(t *testing.T) {

	temp, _ := ioutil.TempDir(".", "es")
	defer func() {
		os.RemoveAll(temp)
	}()

	es := eventstore.NewEventStore(temp)
	ImageCacheEvents(es.RegisterEvent)

	ic, err := loadCache(es)
	assert.NoError(t, err)
	assert.Equal(t, cacheID, ic.AggregateID())

	fromStore := blankImageCache()
	assert.NoError(t, es.LoadAggregate(cacheID, fromStore))
}

func TestCacheAlreadyExists(t *testing.T) {

	temp, _ := ioutil.TempDir(".", "es")
	defer func() {
		os.RemoveAll(temp)
	}()

	es := eventstore.NewEventStore(temp)
	ImageCacheEvents(es.RegisterEvent)

	_, err := loadCache(es)
	assert.NoError(t, err)

	ic, err := loadCache(es)
	assert.NoError(t, err)
	assert.Equal(t, cacheID, ic.AggregateID())
}
