package background

import (
	"brickrecon/eventstore"
	"brickrecon/eventstore/backend/fs"
	"context"
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCacheCreation(t *testing.T) {

	ctx := context.TODO()
	temp, _ := ioutil.TempDir(".", "m")
	defer func() {
		os.RemoveAll(temp)
	}()

	storePath := path.Join(temp, "img/parts")
	os.MkdirAll(storePath, os.ModePerm)

	be, _ := fs.NewFileSystemBackend(path.Join(temp, "es"))
	es := eventstore.NewEventStore(be)
	ImageCacheEvents(ctx, es.RegisterEvent)

	ic, err := NewImageCache(es, storePath, ctx)
	assert.NoError(t, err)
	assert.Equal(t, cacheID, ic.AggregateID())

	fromStore := blankImageCache(storePath)
	assert.NoError(t, es.LoadAggregate(ctx, cacheID, fromStore))
}

func TestCacheAlreadyExists(t *testing.T) {

	ctx := context.TODO()
	temp, _ := ioutil.TempDir(".", "es")
	defer func() {
		os.RemoveAll(temp)
	}()
	storePath := path.Join(temp, "img/parts")
	os.MkdirAll(storePath, os.ModePerm)

	be, _ := fs.NewFileSystemBackend(temp)
	es := eventstore.NewEventStore(be)
	ImageCacheEvents(ctx, es.RegisterEvent)

	_, err := NewImageCache(es, storePath, ctx)
	assert.NoError(t, err)

	ic, err := NewImageCache(es, storePath, ctx)
	assert.NoError(t, err)
	assert.Equal(t, cacheID, ic.AggregateID())
}
