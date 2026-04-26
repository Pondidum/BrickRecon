package command

import (
	"brickrecon/util"
	"fmt"
	"io"
	"reflect"
	"strings"
)

type Renderer = func(w io.Writer, thing any) error

func Render(rendererType string, w io.Writer, thing any) error {
	return TableRenderer(w, thing)
}

func TableRenderer(w io.Writer, thing any) error {

	t := reflect.ValueOf(thing)
	if t.Kind() != reflect.Slice {
		return nil
	}

	lines := make([]string, 0, t.Len()+1)
	fields := reflect.VisibleFields(reflect.TypeOf(thing).Elem())

	th := strings.Builder{}
	for i, h := range fields {
		if i > 0 {
			th.WriteString(" | ")
		}
		th.WriteString(h.Name)
	}
	lines = append(lines, th.String())

	for i := 0; i < t.Len(); i++ {
		elem := t.Index(i)

		values := make([]string, 0, len(fields))
		for _, field := range fields {
			values = append(values, fmt.Sprint(elem.FieldByName(field.Name)))
		}
		lines = append(lines, strings.Join(values, " | "))
	}

	fmt.Fprintln(w, util.TableOutput(lines))

	return nil
}
