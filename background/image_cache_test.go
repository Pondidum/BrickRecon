package background

import (
	"brickrecon/eventstore"
	"brickrecon/lego"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadingDiskCache(t *testing.T) {

	temp, _ := ioutil.TempDir(".", "es")
	defer func() {
		os.RemoveAll(temp)
	}()

	ic := blankImageCache(temp)

	assert.NoError(t, ic.writeFile("3024-85.png", []byte("image one")))
	assert.NoError(t, ic.writeFile("3024-11.png", []byte("image two")))

	ic.ReadFromCache()

	assert.True(t, ic.done["3024-85"])
	assert.True(t, ic.done["3024-11"])
}

func TestReadingDiskCacheDoesntEmitEventsForAlreadyProcessedParts(t *testing.T) {

	temp, _ := ioutil.TempDir(".", "es")
	defer func() {
		os.RemoveAll(temp)
	}()

	ic := blankImageCache(temp)

	assert.NoError(t, ic.writeFile("3024-85.png", []byte("image one")))
	assert.NoError(t, ic.writeFile("3024-11.png", []byte("image two")))

	ic.onFinished(lego.LDrawPart("3024"), 85)
	ic.onFinished(lego.LDrawPart("3024"), 11)

	ic.ReadFromCache()

	assert.Empty(t, eventstore.ReadChanges(ic))
}

func TestReadingDiskCacheDoesntEmitEventsForAlreadyPendingParts(t *testing.T) {

	temp, _ := ioutil.TempDir(".", "es")
	defer func() {
		os.RemoveAll(temp)
	}()

	ic := blankImageCache(temp)

	assert.NoError(t, ic.writeFile("3024-85.png", []byte("image one")))
	assert.NoError(t, ic.writeFile("3024-11.png", []byte("image two")))

	ic.pending["3024-85"] = lego.Part{}
	ic.pending["3024-11"] = lego.Part{}

	ic.ReadFromCache()

	assert.Empty(t, eventstore.ReadChanges(ic))
}
