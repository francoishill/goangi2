package emailUtils

import (
	"fmt"
	"github.com/astaxie/beego"
	"html/template"
)

func renderEmailLink(link *EmailLink) template.HTML {
	if link == nil {
		return template.HTML("")
	}

	html := fmt.Sprintf(
		`<span>%s</span> <a style="%s" target="%s" href="%s">%s</a> <span>%s</span>`,
		link.BeforeLinkText, link.StyleAttributeContent, link.TargetAttribute, link.Href, link.LinkText, link.AfterLinkText)

	return template.HTML(html)
}

func init() {
	beego.AddFuncMap("renderemaillink", renderEmailLink)
}
