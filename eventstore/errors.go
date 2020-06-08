package eventstore

import (
	"fmt"

	uuid "github.com/satori/go.uuid"
)

type AggregateNotFoundError struct {
	ID uuid.UUID
}

func (e *AggregateNotFoundError) Error() string {
	return fmt.Sprintf("No aggregate found for ID %s", e.ID)
}

func IsAggregateNotFound(err error) bool {
	_, ok := err.(*AggregateNotFoundError)

	return ok
}
