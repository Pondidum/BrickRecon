package background

import (
	"brickrecon/eventstore"
	"brickrecon/lego"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"path"
	"strconv"
	"strings"

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

	PartID   lego.PartID
	ColourID int
}

type PartAttempted struct {
	eventstore.EventMeta

	PartID   lego.PartID
	ColourID int
	Error    string
}

type PartImageNotFound struct {
	eventstore.EventMeta

	PartID   lego.PartID
	ColourID int
}

type PartImageStored struct {
	eventstore.EventMeta

	PartID   lego.PartID
	ColourID int
}

type PartAddedFromCache struct {
	eventstore.EventMeta

	PartID   lego.PartID
	ColourID int
}

func ImageCacheEvents(ctx context.Context, register func(context.Context, eventstore.Initialiser)) {
	register(ctx, func() interface{} { return &ImageCacheCreated{} })
	register(ctx, func() interface{} { return &PartAddedFromCache{} })
	register(ctx, func() interface{} { return &PartImageRequested{} })
	register(ctx, func() interface{} { return &PartFetchAttemptsExceeded{} })
	register(ctx, func() interface{} { return &PartAttempted{} })
	register(ctx, func() interface{} { return &PartImageNotFound{} })
	register(ctx, func() interface{} { return &PartImageStored{} })
}

type ImageCache struct {
	*eventstore.Aggregator

	location string

	done     map[string]bool
	pending  map[string]lego.Part
	attempts map[string]int

	writeFile func(filename string, content []byte) error
	listFiles func() ([]string, error)
}

var cacheID uuid.UUID = uuid.Must(uuid.FromString("b83e7c15-24d7-4f18-8de7-34de416eb9de"))

func NewImageCache(es eventstore.EventStore, location string, context context.Context) (*ImageCache, error) {

	ic := blankImageCache(location)

	err := es.LoadAggregate(context, cacheID, ic)

	if err == nil {
		return ic, nil
	}

	if !eventstore.IsAggregateNotFound(err) {
		return nil, err
	}

	ic = createImageCache(cacheID, location)

	if err := ic.ReadFromCache(); err != nil {
		return nil, err
	}

	if err = es.SaveAggregate(context, ic); err != nil {
		return nil, err
	}

	return ic, nil
}

func blankImageCache(location string) *ImageCache {
	ic := &ImageCache{
		location: location,
		done:     map[string]bool{},
		pending:  map[string]lego.Part{},
		attempts: map[string]int{},

		writeFile: func(filename string, content []byte) error {
			return ioutil.WriteFile(path.Join(location, filename), content, 0666)
		},

		listFiles: func() ([]string, error) {
			entries, err := ioutil.ReadDir(location)
			if err != nil {
				return nil, err
			}

			files := []string{}
			for _, entry := range entries {
				if !entry.IsDir() {
					files = append(files, entry.Name())
				}
			}
			return files, nil
		},
	}

	ic.Aggregator = eventstore.NewAggregator(ic.on)

	return ic
}

func createImageCache(id uuid.UUID, location string) *ImageCache {
	ic := blankImageCache(location)
	ic.Apply(&ImageCacheCreated{ID: id})

	return ic
}

func (ic *ImageCache) ReadFromCache() error {

	files, err := ic.listFiles()
	if err != nil {
		return err
	}

	for _, file := range files {
		name := strings.TrimSuffix(file, path.Ext(file))
		parts := strings.Split(name, "-")

		partID := lego.NewPartID(parts[0])
		colourID, err := strconv.Atoi(parts[1])

		if err != nil {
			continue
		}

		if ic.containsPart(name) {
			continue
		}

		ic.Apply(&PartAddedFromCache{PartID: partID, ColourID: colourID})
	}

	return nil
}

func (ic *ImageCache) AddPart(part lego.Part) {

	key := key(part.ID, part.Colour.ID)

	if ic.containsPart(key) {
		return
	}

	ic.Apply(&PartImageRequested{Part: part})
}

func (ic *ImageCache) containsPart(key string) bool {

	if _, found := ic.pending[key]; found {
		return true
	}

	if _, found := ic.done[key]; found {
		return true
	}

	return false
}

func (ic *ImageCache) Run(ctx context.Context) {

	for key, part := range ic.pending {
		fsm := ic.newImageFsm(part.ID, part.Colour.ID)
		fsm.attempts = ic.attempts[key]

		fsm.Run(ctx)

		for _, e := range fsm.events {
			ic.Apply(e)
		}

	}
}

func (ic *ImageCache) on(event eventstore.Event) {

	switch e := event.(type) {

	case *ImageCacheCreated:
		ic.SetID(e.ID)

	case *PartAddedFromCache:
		ic.onFinished(e.PartID, e.ColourID)

	case *PartImageRequested:
		key := key(e.Part.ID, e.Part.Colour.ID)
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

func (ic *ImageCache) onFinished(partID lego.PartID, colourID int) {
	key := key(partID, colourID)

	ic.done[key] = true
	delete(ic.attempts, key)
	delete(ic.pending, key)
}

func (ic *ImageCache) newImageFsm(partID lego.PartID, colourID int) *fsm {
	return &fsm{
		partID:      partID,
		colourID:    colourID,
		maxAttempts: 5,

		invalidImage: []byte{}, //read a png into this

		httpClient: &http.Client{},
		writeFile:  ic.writeFile,
	}
}

func key(id lego.PartID, colourID int) string {
	return fmt.Sprintf("%s-%v", id, colourID)
}
