package emailUtils

type EmailLink struct {
	BeforeLinkText        string
	LinkText              string
	AfterLinkText         string
	Href                  string
	TargetAttribute       string
	StyleAttributeContent string
}

func CreateEmailLink(beforeLinkText, linkText, afterLinkText, href, targetAttribute, styleAttributeContent string) *EmailLink {
	return &EmailLink{
		BeforeLinkText:        beforeLinkText,
		LinkText:              linkText,
		AfterLinkText:         afterLinkText,
		Href:                  href,
		TargetAttribute:       targetAttribute,
		StyleAttributeContent: styleAttributeContent,
	}
}

func (this *EmailLink) GetBeforeLinkText() string {
	return this.BeforeLinkText
}

func (this *EmailLink) GetLinkText() string {
	return this.LinkText
}

func (this *EmailLink) GetAfterLinkText() string {
	return this.AfterLinkText
}

func (this *EmailLink) GetHref() string {
	return this.Href
}

func (this *EmailLink) GetTargetAttribute() string {
	return this.TargetAttribute
}

func (this *EmailLink) GetStyleAttributeContent() string {
	return this.StyleAttributeContent
}
