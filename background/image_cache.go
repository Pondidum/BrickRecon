package background

import (
	"context"
	"fmt"
	"io/ioutil"
	"mvc/eventstore"
	"mvc/lego"
	"net/http"

	"github.com/honeycombio/beeline-go"
)

type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type PartImageRequested struct {
	eventstore.Event

	Part *lego.Part
}

type PartFetchAttemptsExceeded struct {
	eventstore.Event

	PartID   string
	ColourID int
}

type PartAttempted struct {
	eventstore.Event

	PartID   string
	ColourID int
	Error    string
}

type PartImageNotFound struct {
	eventstore.Event

	PartID   string
	ColourID int
}

type PartImageStored struct {
	eventstore.Event

	PartID   string
	ColourID int
}

type ImageCache struct {
	*eventstore.Aggregator

	done     map[string]bool
	pending  map[string]*lego.Part
	attempts map[string]int
}

func NewImageCache() *ImageCache {
	ic := &ImageCache{
		done:     map[string]bool{},
		pending:  map[string]*lego.Part{},
		attempts: map[string]int{},
	}

	ic.Aggregator = eventstore.NewAggregator(ic.on)

	return ic
}

func (ic *ImageCache) AddPart(part *lego.Part) {

	key := key(part.LDrawID, part.Colour.LDrawID)

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
		fsm := NewImageFsm(part)
		fsm.attempts = ic.attempts[key]

		fsm.Run()

		for _, e := range fsm.events {
			ic.Apply(e)
		}
	}
}

func (ic *ImageCache) on(event eventstore.Event) {

	switch e := event.(type) {
	case *PartImageRequested:
		key := key(e.Part.LDrawID, e.Part.Colour.LDrawID)
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

func key(id string, colourID int) string {
	return fmt.Sprintf("%s-%v", id, colourID)
}

// --- fsm

type fsm struct {
	partID      string
	colourID    int
	attempts    int
	maxAttempts int

	httpClient HttpClient
	writeFile  func(filename string, content []byte) error
	events     []eventstore.Event

	invalidImage []byte

	ctx         context.Context
	transitions []string
}

func NewImageFsm(part *lego.Part) *fsm {
	return &fsm{
		partID:      part.LDrawID,
		colourID:    part.Colour.LDrawID,
		maxAttempts: 5,

		invalidImage: []byte{}, //read a png into this

		httpClient: &http.Client{},
		writeFile: func(filename string, content []byte) error {
			return ioutil.WriteFile(filename, content, 0666)
		},
	}
}

func (s *fsm) Run() {

	state := ProcessPart

	ctx, span := beeline.StartSpan(context.Background(), "image_cache")
	defer span.Send()

	s.ctx = ctx

	beeline.AddField(ctx, "part_id", s.partID)
	beeline.AddField(ctx, "colour_id", s.colourID)
	beeline.AddField(ctx, "attempts", s.attempts)
	beeline.AddField(ctx, "max_attempts", s.maxAttempts)

	for state != nil {
		state = state(s)
	}
}

func (s *fsm) httpGet(url string) (*http.Response, error) {

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	return s.httpClient.Do(req)
}

func (s *fsm) event(e eventstore.Event) {
	s.events = append(s.events, e)
}

func (s *fsm) enter(name string) context.Context {
	s.transitions = append(s.transitions, name)
	c, _ := beeline.StartSpan(s.ctx, name)

	return c
}

type state func(s *fsm) (next state)

func ProcessPart(s *fsm) state {

	c := s.enter("process_part")
	maxExceeded := s.attempts < s.maxAttempts

	beeline.AddField(c, "max_attempts_exceeded", maxExceeded)

	if maxExceeded {
		return fetchPart
	}

	return partFailed
}

func partFailed(s *fsm) state {
	s.enter("part_failed")

	s.event(&PartFetchAttemptsExceeded{
		PartID:   s.partID,
		ColourID: s.colourID,
	})

	return storePart(s.invalidImage)
}

func fetchPart(s *fsm) state {
	s.enter("fetch_part")

	url := fmt.Sprintf(`https://img.bricklink.com/ItemImage/PN/%v/%s.png`, s.colourID, s.partID)
	res, err := s.httpGet(url)

	if err != nil {
		return fetchFailed(err)
	}

	defer res.Body.Close()

	if res.StatusCode == http.StatusNotFound {
		return invalidPart
	}

	if res.StatusCode < 200 || res.StatusCode > 299 {
		return fetchFailed(fmt.Errorf("Unexpected statusCode: %v", res.StatusCode))
	}

	content, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return fetchFailed(err)
	}

	return storePart(content)
}

func fetchFailed(err error) state {
	return func(s *fsm) state {
		s.enter("fetch_failed")

		s.event(&PartAttempted{
			PartID:   s.partID,
			ColourID: s.colourID,
			Error:    err.Error(),
		})

		return nil
	}
}

func invalidPart(s *fsm) state {
	s.enter("invalid_part")

	s.event(&PartImageNotFound{
		PartID:   s.partID,
		ColourID: s.colourID,
	})

	return storePart(s.invalidImage)
}

func storePart(content []byte) state {
	return func(s *fsm) state {
		s.enter("store_part")

		filename := fmt.Sprintf("%s-%v.png", s.partID, s.colourID)

		if err := s.writeFile(filename, content); err != nil {
			return storeFailed(err)
		}

		s.event(&PartImageStored{
			PartID:   s.partID,
			ColourID: s.colourID,
		})

		return nil
	}

}

func storeFailed(err error) state {
	return func(s *fsm) state {
		s.enter("store_failed")

		s.event(&PartAttempted{
			PartID:   s.partID,
			ColourID: s.colourID,
			Error:    err.Error(),
		})
		return nil
	}
}
