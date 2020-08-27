package preen

import (
	"brickrecon/util"
	"regexp"
	"strings"
)

var rx = regexp.MustCompile("{(.*?)}")

type ControllerLinker func(controller string, parameters map[string]interface{}) string

func NewControllerLinker(controllers []Controller) ControllerLinker {

	lookup := map[string]string{}

	for _, ctl := range controllers {
		lookup[controllerName(ctl)] = ctl.Path()
	}

	return func(controller string, parameters map[string]interface{}) string {

		toControllerPath := lookup[controller]

		url := "/" + rx.ReplaceAllStringFunc(toControllerPath, func(match string) string {
			return findValue(parameters, strings.Trim(match, "{}"))
		})

		return url
	}
}

func findValue(source map[string]interface{}, key string) string {
	for k, v := range source {
		if strings.EqualFold(key, k) {
			return util.Strval(v)
		}
	}
	return ""
}
