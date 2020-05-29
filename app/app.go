package app

import "mvc/lego"

type SiteModel struct {
	AllModels []string

	SelectedModel lego.Model
}
