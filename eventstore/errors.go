package eventstore

import (
	"fmt"
)

type AggregateNotFoundError struct {
	ID AggregateID
}

func (e *AggregateNotFoundError) Error() string {
	return fmt.Sprintf("No aggregate found for ID %s", e.ID)
}

func IsAggregateNotFound(err error) bool {
	_, ok := err.(*AggregateNotFoundError)

	return ok
}
