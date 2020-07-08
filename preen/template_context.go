package preen

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"github.com/gorilla/context"
)

func TemplateFuncDefinitions() template.FuncMap {

	return template.FuncMap{
		"_user": func() UserInfo {
			return UserInfo{}
		},
		"_page": func() PageInfo {
			return PageInfo{}
		},
		"_site": func() SiteInfo {
			return SiteInfo{}
		},
		"active": func(parts ...interface{}) bool {
			return false
		},
		"activeWith": func(url string, queries ...interface{}) bool {
			return false
		},
	}
}

func TemplateFuncs(req *http.Request) template.FuncMap {
	return template.FuncMap{
		"_user": func() UserInfo {
			user, found := context.Get(req, "UserInfo").(UserInfo)

			if !found {
				user = UserInfo{}
			}

			return user
		},
		"_page": func() PageInfo {
			return PageInfo{Path: req.URL.Path}
		},
		"_site": func() SiteInfo {
			return SiteInfo{URL: req.Host}
		},
		"active": func(parts ...interface{}) bool {
			path := ""

			for _, p := range parts {
				path = path + strval(p)
			}

			urlPath := req.URL.Path

			if req.URL.RawQuery != "" {
				urlPath += "?" + req.URL.RawQuery
			}
			return strings.HasPrefix(urlPath, path)
		},
		"activeWith": func(url string, queries ...interface{}) bool {

			if req.URL.Path != url {
				return false
			}

			qs := req.URL.Query()

			for i := 0; i < len(queries); i += 2 {
				key := strval(queries[i])
				value := strval(queries[i+1])

				if qs.Get(key) != value {
					return false
				}
			}

			return true
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
