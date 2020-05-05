package command

import "github.com/posener/complete"

func mergeAutocompleteFlags(flags ...complete.Flags) complete.Flags {
	merged := make(map[string]complete.Predictor, len(flags))

	for _, f := range flags {
		for k, v := range f {
			merged[k] = v
		}
	}

	return merged
}
