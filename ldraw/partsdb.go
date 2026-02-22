package ldraw

import (
	"archive/zip"
	"brickrecon/tracing"
	"bufio"
	"bytes"
	"context"
	"io"
	"path"
	"strings"
	"unicode"

	"go.opentelemetry.io/otel"
)

var tr = otel.Tracer("goes")

type partDto struct {
	Name       string
	PartNumber string
	MovedTo    string

	AlternateIds []string
}

type LDrawPart struct {
	PartNumber     string
	OldPartNumbers []string
}

func ParseDatabaseArchive(ctx context.Context, archive io.Reader) (map[string]*partDto, error) {
	ctx, span := tr.Start(ctx, "parse_database_archive")
	defer span.End()

	content, err := io.ReadAll(archive)
	if err != nil {
		return nil, tracing.Error(span, err)
	}

	zr, err := zip.NewReader(bytes.NewReader(content), int64(len(content)))
	if err != nil {
		return nil, tracing.Error(span, err)
	}

	parts := map[string]*partDto{}
	for _, file := range zr.File {

		dto, err := parsePart(file)
		if err != nil {
			return nil, tracing.Error(span, err)
		}
		if dto == nil {
			continue
		}
		// dtos = append(dtos, dto)
		parts[dto.PartNumber] = dto
	}

	calculateNewestMoves(parts)

	return parts, nil
}

func parsePart(file *zip.File) (*partDto, error) {

	if !strings.HasPrefix(file.Name, "ldraw/parts/") {
		return nil, nil
	}
	if !strings.HasSuffix(file.Name, ".dat") {
		return nil, nil
	}

	partNum := strings.TrimSuffix(path.Base(file.Name), ".dat")
	// we are ignoring `t` (thirdparty) and `u` (unofficial) parts for now
	if !unicode.IsDigit(rune(partNum[0])) {
		return nil, nil
	}

	fr, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer fr.Close()

	part := parseDat(fr)
	part.PartNumber = partNum

	return part, nil
}

func parseDat(fr io.Reader) *partDto {

	part := &partDto{}
	scanner := bufio.NewScanner(fr)

	first := true
	for scanner.Scan() {
		if len(scanner.Bytes()) == 0 {
			continue
		}

		if len(bytes.TrimSpace(scanner.Bytes())) == 0 {
			continue
		}

		line := scanner.Text()

		if target, ok := strings.CutPrefix(line, "0 ~Moved to "); ok {
			part.MovedTo = target
		} else if first {

			if name, ok := strings.CutPrefix(scanner.Text(), "0 "); ok {
				part.Name = name
			}
		}

		first = false

		// https://www.ldraw.org/article/340.html#keywords
		if raw, ok := strings.CutPrefix(line, "0 !KEYWORDS "); ok {

			for statement := range strings.SplitSeq(raw, ", ") {
				if id, ok := strings.CutPrefix(statement, "BrickLink "); ok {
					part.AlternateIds = append(part.AlternateIds, strings.TrimSpace(id))
				} else if id, ok := strings.CutPrefix(statement, "Rebrickable "); ok {
					part.AlternateIds = append(part.AlternateIds, strings.TrimSpace(id))
				}
			}
		}

		if line[0] != '0' {
			break
		}
	}

	return part
}

func calculateNewestMoves(input map[string]*partDto) {

	for _, part := range input {
		if part.MovedTo == "" {
			continue
		}

		for {
			newerPart, found := input[part.MovedTo]
			if !found || newerPart.MovedTo == "" {
				break
			}
			part.MovedTo = newerPart.MovedTo
		}
	}
}
