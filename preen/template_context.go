package preen

import (
	"html/template"
	"net/http"

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
	}
}

type SiteInfo struct {
	URL string
}

type PageInfo struct {
	Path string
}
