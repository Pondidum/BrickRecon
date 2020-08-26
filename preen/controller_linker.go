package preen

import (
	"regexp"
	"strings"
)

var rx = regexp.MustCompile("{(.*?)}")

type ControllerLinker func(controller string, parameters map[string]string) string

func CreateControllerLinker(controllers []Controller) ControllerLinker {

	lookup := map[string]string{}

	for _, ctl := range controllers {
		lookup[controllerName(ctl)] = ctl.Path()
	}

	return func(controller string, parameters map[string]string) string {

		toControllerPath := lookup[controller]

		url := "/" + rx.ReplaceAllStringFunc(toControllerPath, func(match string) string {
			return parameters[strings.Trim(match, "{}")]
		})

		return url
	}
}
