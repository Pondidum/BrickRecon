package command

import (
	"brickrecon/util"
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"strings"
)

type Renderer = func(w io.Writer, thing any) error

func Render(rendererType string, w io.Writer, thing any) error {
	switch strings.ToLower(rendererType) {
	case "table":
		return TableRenderer(w, thing)
	case "json":
		return JsonRenderer(w, thing)
	default:
		return fmt.Errorf("unsupported renderer: %s", rendererType)
	}
}

func TableRenderer(w io.Writer, thing any) error {

	t := reflect.ValueOf(thing)
	if t.Kind() != reflect.Slice {
		return nil
	}

	lines := make([]string, 0, t.Len()+1)
	headers := getHeaders(thing)

	th := strings.Builder{}
	for i, h := range headers {
		if i > 0 {
			th.WriteString(" | ")
		}
		th.WriteString(h)
	}
	lines = append(lines, th.String())

	for i := 0; i < t.Len(); i++ {
		elem := t.Index(i)

		values := make([]string, 0, len(headers))
		for _, header := range headers {
			values = append(values, fmt.Sprint(elem.FieldByName(header)))
		}
		lines = append(lines, strings.Join(values, " | "))
	}

	fmt.Fprintln(w, util.TableOutput(lines))

	return nil
}

func getHeaders(thing any) []string {

	fields := reflect.VisibleFields(reflect.TypeOf(thing).Elem())
	headers := make([]string, 0, len(fields))

	for _, field := range fields {
		if field.IsExported() {
			headers = append(headers, field.Name)
		}
	}

	return headers
}

func JsonRenderer(w io.Writer, thing any) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(thing)
}
