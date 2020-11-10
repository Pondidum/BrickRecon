package background

import (
	"brickrecon/eventstore"
	"brickrecon/lego"
	"context"
	"io/ioutil"
	"net/http"
	"path"
	"strconv"
	"strings"
)

type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type ImageCacheCreated struct {
	eventstore.EventMeta

	ID eventstore.AggregateID
}

type PartImageRequested struct {
	eventstore.EventMeta

	Part *lego.Part
}

type PartFetchAttemptsExceeded struct {
	eventstore.EventMeta

	PartID   lego.LDrawPart
	ColourID lego.LDrawColour
}

type PartAttempted struct {
	eventstore.EventMeta

	PartID   lego.LDrawPart
	ColourID lego.LDrawColour
	Error    string
}

type PartImageNotFound struct {
	eventstore.EventMeta

	PartID   lego.LDrawPart
	ColourID lego.LDrawColour
}

type PartImageStored struct {
	eventstore.EventMeta

	PartID   lego.LDrawPart
	ColourID lego.LDrawColour
}

type PartAddedFromCache struct {
	eventstore.EventMeta

	PartID   lego.LDrawPart
	ColourID lego.LDrawColour
}

var ImageCacheEvents = []eventstore.Initialiser{
	func() interface{} { return &ImageCacheCreated{} },
	func() interface{} { return &PartAddedFromCache{} },
	func() interface{} { return &PartImageRequested{} },
	func() interface{} { return &PartFetchAttemptsExceeded{} },
	func() interface{} { return &PartAttempted{} },
	func() interface{} { return &PartImageNotFound{} },
	func() interface{} { return &PartImageStored{} },
}

type ImageCache struct {
	*eventstore.Aggregator

	location string

	done     map[lego.PartKey]bool
	pending  map[lego.PartKey]*lego.Part
	attempts map[lego.PartKey]int

	writeFile func(filename string, content []byte) error
	listFiles func() ([]string, error)
}

var cacheID = eventstore.AggregateID("b83e7c15-24d7-4f18-8de7-34de416eb9de")

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
		done:     map[lego.PartKey]bool{},
		pending:  map[lego.PartKey]*lego.Part{},
		attempts: map[lego.PartKey]int{},

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

func createImageCache(id eventstore.AggregateID, location string) *ImageCache {
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

		partID := lego.LDrawPart(parts[0])
		colourID, err := strconv.Atoi(parts[1])

		if err != nil {
			continue
		}

		if ic.containsPart(lego.CreatePartKey(partID, lego.LDrawColour(colourID))) {
			continue
		}

		ic.Apply(&PartAddedFromCache{
			PartID:   partID,
			ColourID: lego.LDrawColour(colourID),
		})
	}

	return nil
}

func (ic *ImageCache) AddPart(part *lego.Part) {

	if ic.containsPart(part.Key) {
		return
	}

	ic.Apply(&PartImageRequested{Part: part})
}

func (ic *ImageCache) containsPart(key lego.PartKey) bool {

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
		fsm := ic.newImageFsm(part.Aliases.LDrawID, part.Colour.Aliases.LDrawID)
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
		ic.pending[e.Part.Key] = e.Part

	case *PartAttempted:
		key := lego.CreatePartKey(e.PartID, e.ColourID)
		ic.attempts[key]++

	case *PartFetchAttemptsExceeded:
		ic.onFinished(e.PartID, e.ColourID)

	case *PartImageNotFound:
		ic.onFinished(e.PartID, e.ColourID)

	case *PartImageStored:
		ic.onFinished(e.PartID, e.ColourID)

	}
}

func (ic *ImageCache) onFinished(partID lego.LDrawPart, colourID lego.LDrawColour) {
	key := lego.CreatePartKey(partID, colourID)

	ic.done[key] = true
	delete(ic.attempts, key)
	delete(ic.pending, key)
}

func (ic *ImageCache) newImageFsm(partID lego.LDrawPart, colourID lego.LDrawColour) *fsm {
	return &fsm{
		partID:      partID,
		colourID:    colourID,
		maxAttempts: 5,

		httpClient: &http.Client{},
		writeFile:  ic.writeFile,
	}
}
