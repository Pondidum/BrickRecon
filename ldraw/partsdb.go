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
	PartNumber string
	MovedTo    string
}

type LDrawPart struct {
	PartNumber     string
	OldPartNumbers []string
}

func ParseDatabaseArchive(ctx context.Context, archive io.Reader) (map[string]string, error) {
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

	dtos := map[string]string{}
	for _, file := range zr.File {

		dto, err := parsePart(file)
		if err != nil {
			return nil, tracing.Error(span, err)
		}
		if dto == nil {
			continue
		}
		// dtos = append(dtos, dto)
		dtos[dto.PartNumber] = dto.MovedTo
	}

	parts := buildParts(dtos)

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

	return &partDto{
		PartNumber: partNum,
		MovedTo:    movedTo(fr),
	}, nil
}

func movedTo(fr io.Reader) string {

	scanner := bufio.NewScanner(fr)
	for scanner.Scan() {
		if target, ok := strings.CutPrefix(scanner.Text(), "0 ~Moved to "); ok {
			return target
		}
	}
	return ""
}

func buildParts(input map[string]string) map[string]string {

	for part, movedTo := range input {
		if movedTo == "" {
			continue
		}

		for {
			newer, found := input[movedTo]
			if !found || newer == "" {
				break
			}
			movedTo = newer
		}

		input[part] = movedTo
	}

	return input
}
