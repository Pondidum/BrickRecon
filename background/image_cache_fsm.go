package background

import (
	"context"
	"fmt"
	"io/ioutil"
	"mvc/eventstore"
	"net/http"

	"github.com/honeycombio/beeline-go"
)

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

type state func(s *fsm) (next state)

func (s *fsm) Run() {

	state := processPart

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

func processPart(s *fsm) state {

	c := s.enter("process_part")
	maxExceeded := s.attempts >= s.maxAttempts

	beeline.AddField(c, "max_attempts_exceeded", maxExceeded)

	if maxExceeded {
		return partFailed
	}

	return fetchPart
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
	c := s.enter("fetch_part")

	url := fmt.Sprintf(`https://img.bricklink.com/ItemImage/PN/%v/%s.png`, s.colourID, s.partID)
	res, err := s.httpGet(url)

	if err != nil {
		beeline.AddField(c, "error", err)
		return fetchFailed(err)
	}

	defer res.Body.Close()

	beeline.AddField(c, "status_code", res.StatusCode)

	if res.StatusCode == http.StatusNotFound {
		return invalidPart
	}

	if res.StatusCode < 200 || res.StatusCode > 299 {
		return fetchFailed(fmt.Errorf("Unexpected statusCode: %v", res.StatusCode))
	}

	content, err := ioutil.ReadAll(res.Body)

	if err != nil {
		beeline.AddField(c, "error", err)

		return fetchFailed(err)
	}

	beeline.AddField(c, "content_length", len(content))

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
