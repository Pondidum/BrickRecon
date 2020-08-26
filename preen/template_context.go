package preen

import (
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/gorilla/context"
)

type TemplateContext struct {
	Request *http.Request

	Functions template.FuncMap
}

func NewTemplateContext() *TemplateContext {
	tf := &TemplateContext{}
	tf.Functions = realFunctions(tf)

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
				key := strval(queries[i])
				value := strval(queries[i+1])

				if qs.Get(key) != value {
					return false
				}
			}

			return true
		},
		"format": func(ts time.Time, layout string) string {
			return ts.Format(layout)
		},
	}
}

func strval(v interface{}) string {
	switch v := v.(type) {
	case string:
		return v
	case []byte:
		return string(v)
	case error:
		return v.Error()
	case fmt.Stringer:
		return v.String()
	default:
		return fmt.Sprintf("%v", v)
	}
}

type SiteInfo struct {
	URL string
}

type PageInfo struct {
	Path string
}
