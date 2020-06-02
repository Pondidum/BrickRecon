package lego

import (
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"strings"
)

const (
	brickLinkID int = iota
	elementID
	ldrawID
	partName
	brickLinkColour
	ldrawColour
	colourName
	colourCategory
	quantity
	weight
)

func ReadPartsList(content io.Reader) ([]Part, error) {

	reader := csv.NewReader(content)
	reader.Comma = '\t'

	parts := []Part{}

	// read the header
	_, err := reader.Read()
	if err != nil {
		return nil, err
	}

	for {
		record, err := reader.Read()

		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, err
		}

		if isBlank(record) {
			continue
		}

		if isSummaryHeader(record) {
			break
		}

		part, err := parsePart(record)

		if err != nil {
			return nil, fmt.Errorf("%s\nRecord: %s", err.Error(), record)
		}

		parts = append(parts, part)
	}

	return parts, nil
}

func isBlank(fields []string) bool {
	return len(fields) == 0 || strings.TrimSpace(fields[0]) == ""
}

func isSummaryHeader(fields []string) bool {
	return len(fields) > 0 && fields[brickLinkID] == "Total qty"
}

func parsePart(fields []string) (Part, error) {

	var err error
	part := Part{
		BrickLinkID: fields[brickLinkID],
		ElementID:   parseElementID(fields[elementID]),
		LDrawID:     fields[ldrawID],
		Name:        fields[partName],
	}

	if part.Colour, err = parseColour(fields); err != nil {
		return Part{}, err
	}

	if part.Quantity, err = strconv.Atoi(fields[quantity]); err != nil {
		return Part{}, convertError("part.Quantity", fields[quantity])
	}

	if part.Weight, err = strconv.ParseFloat(fields[weight], 64); err != nil {
		return Part{}, convertError("part.Weight", fields[weight])
	}

	return part, err
}

func parseElementID(value string) int {

	if value == "" {
		return 0
	}

	id, err := strconv.Atoi(value)

	if err != nil {
		return 0
	}

	return id
}

func parseColour(fields []string) (Colour, error) {

	var err error
	colour := Colour{
		Name:     fields[colourName],
		Category: fields[colourCategory],
	}

	if colour.BrickLinkID, err = strconv.Atoi(fields[brickLinkColour]); err != nil {
		return Colour{}, convertError("colour.BrickLinkID", fields[brickLinkID])
	}

	if colour.LDrawID, err = strconv.Atoi(fields[ldrawColour]); err != nil {
		return Colour{}, convertError("colour.LDrawID", fields[ldrawColour])
	}

	return colour, err
}

func convertError(key string, value string) error {
	return fmt.Errorf("Unable to convert '%s' to %s", value, key)
}
