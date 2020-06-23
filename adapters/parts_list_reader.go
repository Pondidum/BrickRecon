package adapters

import (
	"brickrecon/lego"
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

func ReadPartsList(content io.Reader) ([]lego.Part, error) {

	reader := csv.NewReader(content)
	reader.Comma = '\t'

	parts := []lego.Part{}

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

func parsePart(fields []string) (lego.Part, error) {

	var err error

	part := lego.Part{
		Name:    fields[partName],
		Aliases: parsePartAliases(fields),
		ID:      lego.NewPartID(fields[ldrawID]),
	}

	if part.Colour, err = parseColour(fields); err != nil {
		return lego.Part{}, err
	}

	if part.Quantity, err = strconv.Atoi(fields[quantity]); err != nil {
		return lego.Part{}, convertError("part.Quantity", fields[quantity])
	}

	if part.Weight, err = strconv.ParseFloat(fields[weight], 64); err != nil {
		return lego.Part{}, convertError("part.Weight", fields[weight])
	}

	return part, err
}

func parsePartAliases(fields []string) lego.PartAliases {
	return lego.PartAliases{
		BrickLinkID: fields[brickLinkID],
		ElementID:   parseElementID(fields[elementID]),
		LDrawID:     fields[ldrawID],
	}

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

func parseColour(fields []string) (lego.Colour, error) {

	var err error
	var bricklinkID, ldrawID int

	if bricklinkID, err = strconv.Atoi(fields[brickLinkColour]); err != nil {
		return lego.Colour{}, convertError("colour.BrickLinkID", fields[brickLinkID])
	}

	if ldrawID, err = strconv.Atoi(fields[ldrawColour]); err != nil {
		return lego.Colour{}, convertError("colour.LDrawID", fields[ldrawColour])
	}

	aliases := lego.ColourAliases{
		BrickLinkID: lego.BrickLinkColour(bricklinkID),
		LDrawID:     lego.LDrawColour(ldrawID),
	}

	colour := lego.Colour{
		ID:       aliases.LDrawID,
		Aliases:  aliases,
		Name:     fields[colourName],
		Category: fields[colourCategory],
	}

	return colour, err
}

func convertError(key string, value string) error {
	return fmt.Errorf("Unable to convert '%s' to %s", value, key)
}
