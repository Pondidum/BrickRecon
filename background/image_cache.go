package background

import (
	"context"
	"fmt"
	"io/ioutil"
	"mvc/eventstore"
	"net/http"

	"github.com/honeycombio/beeline-go"
)

type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
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

// --- fsm

func Run(s *dto) {

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

type dto struct {
	partID      string
	colourID    int
	attempts    int
	maxAttempts int

	content []byte
	format  string

	httpClient HttpClient
	writeFile  func(filename string, content []byte) error
	events     []eventstore.Event

	invalidImage []byte

	ctx         context.Context
	transitions []string
}

func (s *dto) httpGet(url string) (*http.Response, error) {

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	return s.httpClient.Do(req)
}

func (s *dto) event(e eventstore.Event) {
	s.events = append(s.events, e)
}

func (s *dto) enter(name string) context.Context {
	s.transitions = append(s.transitions, name)
	c, _ := beeline.StartSpan(s.ctx, name)

	return c
}

type state func(s *dto) (next state)

func ProcessPart(s *dto) state {

	c := s.enter("process_part")
	maxExceeded := s.attempts < s.maxAttempts

	beeline.AddField(c, "max_attempts_exceeded", maxExceeded)

	if maxExceeded {
		return fetchPart
	}

	return partFailed
}

func partFailed(s *dto) state {
	s.enter("part_failed")

	s.event(PartFetchAttemptsExceeded{
		PartID:   s.partID,
		ColourID: s.colourID,
	})

	return nil
}

func fetchPart(s *dto) state {
	s.enter("fetch_part")

	url := fmt.Sprintf(`https://img.bricklink.com/ItemImage/PN/%v/%s.%s`, s.colourID, s.partID, s.format)
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

	s.content = content

	return storePart
}

func fetchFailed(err error) state {
	return func(s *dto) state {
		s.enter("fetch_failed")

		s.event(PartAttempted{
			PartID:   s.partID,
			ColourID: s.colourID,
			Error:    err.Error(),
		})

		return nil
	}
}

func invalidPart(s *dto) state {
	s.enter("invalid_part")

	s.event(PartImageNotFound{
		PartID:   s.partID,
		ColourID: s.colourID,
	})

	s.content = s.invalidImage

	return storePart
}

func storePart(s *dto) state {
	s.enter("store_part")

	filename := fmt.Sprintf("%s-%v.%s", s.partID, s.colourID, s.format)

	if err := s.writeFile(filename, s.content); err != nil {
		return storeFailed(err)
	}

	s.event(PartImageStored{
		PartID:   s.partID,
		ColourID: s.colourID,
	})

	return nil
}

func storeFailed(err error) state {
	return func(s *dto) state {
		s.enter("store_failed")

		s.event(PartAttempted{
			PartID:   s.partID,
			ColourID: s.colourID,
			Error:    err.Error(),
		})
		return nil
	}
}
