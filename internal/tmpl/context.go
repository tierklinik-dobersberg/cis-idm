package tmpl

import "strings"

type Context interface {
	SetSiteName(name string)
	SetSiteURL(name string)
	SetPublicURL(url string)
}

type BaseContext struct {
	SiteName  string
	SiteURL   string
	PublicURL string
}

func (bc *BaseContext) SetSiteName(name string) {
	bc.SiteName = name
}

func (bc *BaseContext) SetSiteURL(url string) {
	bc.SiteURL = url
}

func (bc *BaseContext) SetPublicURL(url string) {
	bc.PublicURL = strings.TrimSuffix(url, "/")
}
