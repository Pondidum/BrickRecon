package preen

import (
	"brickrecon/util"
	"html/template"
	"net/http"
	"time"

	"github.com/gorilla/context"
)

type TemplateContext struct {
	Request    *http.Request
	Controller Controller
	Functions  template.FuncMap

	linker ControllerLinker
}

func NewTemplateContext(linker ControllerLinker) *TemplateContext {
	tf := &TemplateContext{}
	tf.Functions = realFunctions(tf)
	tf.linker = linker

	return tf
}

func realFunctions(tf *TemplateContext) template.FuncMap {
	return template.FuncMap{
		"_user": func() UserInfo {
			user, found := context.Get(tf.Request, "UserInfo").(UserInfo)

			if !found {
				user = UserInfo{}
			}

			return user
		},
		"_page": func() PageInfo {
			return PageInfo{Path: tf.Request.URL.Path}
		},
		"_site": func() SiteInfo {
			return SiteInfo{URL: tf.Request.Host}
		},
		"active": func(url string, queries ...interface{}) bool {

			if tf.Request.URL.Path != url {
				return false
			}

			qs := tf.Request.URL.Query()

			for i := 0; i < len(queries); i += 2 {
				key := util.Strval(queries[i])
				value := util.Strval(queries[i+1])

				if qs.Get(key) != value {
					return false
				}
			}

			return true
		},
		"format": func(ts time.Time, layout string) string {
			return ts.Format(layout)
		},
		"linkto": func(controller string, parameters ...interface{}) string {

			args := map[string]interface{}{}

			for i := 0; i < len(parameters); i += 2 {
				key := util.Strval(parameters[i])
				value := parameters[i+1]

				args[key] = value
			}

			return tf.linker(controller, args)
		},
		"html": func(html string) template.HTML {
			return template.HTML(html)
		},
	}
}

type SiteInfo struct {
	URL string
}

type PageInfo struct {
	Path string
}
