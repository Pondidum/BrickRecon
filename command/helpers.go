package command

import (
	"strings"

	"github.com/posener/complete"
	"github.com/ryanuber/columnize"
)

func mergeAutocompleteFlags(flags ...complete.Flags) complete.Flags {
	merged := make(map[string]complete.Predictor, len(flags))

	for _, f := range flags {
		for k, v := range f {
			merged[k] = v
		}
	}

	return merged
}

func tableOutput(list []string) string {
	if len(list) == 0 {
		return ""
	}

	delim := "|"
	underline := ""
	headers := strings.Split(list[0], delim)
	for i, h := range headers {
		h = strings.TrimSpace(h)
		u := strings.Repeat("-", len(h))

		underline = underline + u
		if i != len(headers)-1 {
			underline = underline + delim
		}
	}

	list = append(list, "")
	copy(list[2:], list[1:])
	list[1] = underline

	return columnize.Format(list, &columnize.Config{
		Glue: "    ",
	})
}
