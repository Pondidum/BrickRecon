package background

import (
	"brickrecon/eventstore"
	"brickrecon/eventstore/backend/fs"
	"context"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCacheCreation(t *testing.T) {

	ctx := context.TODO()
	temp, _ := ioutil.TempDir(".", "es")
	defer func() {
		os.RemoveAll(temp)
	}()

	be, _ := fs.NewFileSystemBackend(temp)
	es := eventstore.NewEventStore(be)
	ImageCacheEvents(ctx, es.RegisterEvent)

	ic, err := NewImageCache(es, ctx)
	assert.NoError(t, err)
	assert.Equal(t, cacheID, ic.AggregateID())

	fromStore := blankImageCache("./app/static/img/parts")
	assert.NoError(t, es.LoadAggregate(ctx, cacheID, fromStore))
}

func TestCacheAlreadyExists(t *testing.T) {

	ctx := context.TODO()
	temp, _ := ioutil.TempDir(".", "es")
	defer func() {
		os.RemoveAll(temp)
	}()

	be, _ := fs.NewFileSystemBackend(temp)
	es := eventstore.NewEventStore(be)
	ImageCacheEvents(ctx, es.RegisterEvent)

	_, err := NewImageCache(es, ctx)
	assert.NoError(t, err)

	ic, err := NewImageCache(es, ctx)
	assert.NoError(t, err)
	assert.Equal(t, cacheID, ic.AggregateID())
}
