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
				path = path + fmt.Sprintf("%v", p)
			}

			urlPath := req.URL.Path

			if req.URL.RawQuery != "" {
				urlPath += "?" + req.URL.RawQuery
			}
			return strings.HasPrefix(urlPath, path)
		},
	}
}

type SiteInfo struct {
	URL string
}

type PageInfo struct {
	Path string
}
