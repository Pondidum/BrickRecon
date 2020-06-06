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

type PartFailed struct {
	eventstore.Event
}

type PartAttempted struct {
	eventstore.Event
}

type PartImageNotFound struct {
	eventstore.Event
}

type PartImageStored struct {
	eventstore.Event
}

// --- fsm

func Run(s *dto) {

	state := ProcessPart

	ctx, span := beeline.StartSpan(context.Background(), "image_cache")
	defer span.Send()

	s.ctx = ctx

	for state != nil {
		state = state(s)
	}
}

type dto struct {
	partID   string
	colourID int

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

type State func(s *dto) (next State)

func ProcessPart(s *dto) State {

	c := s.enter("process_part")

	beeline.AddField(c, "attempts", s.attempts)
	beeline.AddField(c, "max_attempts", s.maxAttempts)

	if s.attempts < s.maxAttempts {
		return fetchPart
	}

	return partFailed
}

func partFailed(s *dto) State {
	s.enter("part_failed")

	s.event(PartFailed{})
	return nil
}

func fetchPart(s *dto) State {
	s.enter("fetch_part")

	url := fmt.Sprintf(`https://img.bricklink.com/ItemImage/PN/%v/%s.%s`, s.colourID, s.partID, s.format)
	res, err := s.httpGet(url)

	if err != nil {
		return fetchFailed
	}

	defer res.Body.Close()

	if res.StatusCode == http.StatusNotFound {
		return invalidPart
	}

	if res.StatusCode < 200 || res.StatusCode > 299 {
		return fetchFailed
	}

	content, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return fetchFailed
	}

	s.content = content

	return storePart
}

func fetchFailed(s *dto) State {
	s.enter("fetch_failed")

	s.event(PartAttempted{})
	return nil
}

func invalidPart(s *dto) State {
	s.enter("invalid_part")

	s.event(PartImageNotFound{})
	s.content = s.invalidImage

	return storePart
}

func storePart(s *dto) State {
	s.enter("store_part")

	filename := fmt.Sprintf("%s-%v.%s", s.partID, s.colourID, s.format)

	if err := s.writeFile(filename, s.content); err != nil {
		return storeFailed
	}

	s.event(PartImageStored{})

	return nil
}

func storeFailed(s *dto) State {
	s.enter("store_failed")

	s.event(PartAttempted{})
	return nil
}

// ---
