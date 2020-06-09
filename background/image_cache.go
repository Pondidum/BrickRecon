package background

import (
	"fmt"
	"io/ioutil"
	"mvc/eventstore"
	"mvc/lego"
	"net/http"
	"path"

	uuid "github.com/satori/go.uuid"
)

type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type ImageCacheCreated struct {
	eventstore.EventMeta

	ID uuid.UUID
}

type PartImageRequested struct {
	eventstore.EventMeta

	Part lego.Part
}

type PartFetchAttemptsExceeded struct {
	eventstore.EventMeta

	PartID   string
	ColourID int
}

type PartAttempted struct {
	eventstore.EventMeta

	PartID   string
	ColourID int
	Error    string
}

type PartImageNotFound struct {
	eventstore.EventMeta

	PartID   string
	ColourID int
}

type PartImageStored struct {
	eventstore.EventMeta

	PartID   string
	ColourID int
}

func ImageCacheEvents(register func(eventstore.Initialiser)) {
	register(func() interface{} { return &ImageCacheCreated{} })
	register(func() interface{} { return &PartImageRequested{} })
	register(func() interface{} { return &PartFetchAttemptsExceeded{} })
	register(func() interface{} { return &PartAttempted{} })
	register(func() interface{} { return &PartImageNotFound{} })
	register(func() interface{} { return &PartImageStored{} })
}

type ImageCache struct {
	*eventstore.Aggregator

	location string

	done     map[string]bool
	pending  map[string]lego.Part
	attempts map[string]int
}

func blankImageCache(location string) *ImageCache {
	ic := &ImageCache{
		location: location,
		done:     map[string]bool{},
		pending:  map[string]lego.Part{},
		attempts: map[string]int{},
	}

	ic.Aggregator = eventstore.NewAggregator(ic.on)

	return ic
}

func NewImageCache(id uuid.UUID, location string) *ImageCache {
	ic := blankImageCache(location)
	ic.Apply(&ImageCacheCreated{ID: id})

	return ic
}

func (ic *ImageCache) AddPart(part lego.Part) {

	key := key(part.BrickLinkID, part.Colour.BrickLinkID)

	if _, found := ic.pending[key]; found {
		return
	}

	if _, found := ic.done[key]; found {
		return
	}

	ic.Apply(&PartImageRequested{Part: part})
}

func (ic *ImageCache) Run() {

	for key, part := range ic.pending {
		fsm := ic.newImageFsm(part.BrickLinkID, part.Colour.BrickLinkID)
		fsm.attempts = ic.attempts[key]

		fsm.Run()

		for _, e := range fsm.events {
			ic.Apply(e)
		}
	}
}

func (ic *ImageCache) on(event eventstore.Event) {

	switch e := event.(type) {

	case *ImageCacheCreated:
		ic.SetID(e.ID)

	case *PartImageRequested:
		key := key(e.Part.BrickLinkID, e.Part.Colour.BrickLinkID)
		ic.pending[key] = e.Part

	case *PartAttempted:
		key := key(e.PartID, e.ColourID)
		ic.attempts[key]++

	case *PartFetchAttemptsExceeded:
		ic.onFinished(e.PartID, e.ColourID)

	case *PartImageNotFound:
		ic.onFinished(e.PartID, e.ColourID)

	case *PartImageStored:
		ic.onFinished(e.PartID, e.ColourID)

	}
}

func (ic *ImageCache) onFinished(partID string, colourID int) {
	key := key(partID, colourID)

	ic.done[key] = true
	delete(ic.attempts, key)
	delete(ic.pending, key)
}

func (ic *ImageCache) newImageFsm(partID string, colourID int) *fsm {
	return &fsm{
		partID:      partID,
		colourID:    colourID,
		maxAttempts: 5,

		invalidImage: []byte{}, //read a png into this

		httpClient: &http.Client{},
		writeFile:  ic.writeFile,
	}
}

func key(id string, colourID int) string {
	return fmt.Sprintf("%s-%v", id, colourID)
}
