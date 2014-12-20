package controllers

import (
	"net/url"
	"strings"
)

func (this *BaseController) GetOptionalQueryValueBoolean(queryKey string, defaultVal bool) (hadQueryValue bool, boolVal bool) {
	escapedQueryValue := this.Ctx.Input.Query(queryKey)
	if escapedQueryValue == "" {
		return false, defaultVal
	}
	unescapedQueryValue, err := url.QueryUnescape(escapedQueryValue)
	this.PanicIfError(err)
	unescapedQueryValue = strings.ToLower(strings.Trim(unescapedQueryValue, " "))

	if unescapedQueryValue != "true" && unescapedQueryValue != "1" &&
		unescapedQueryValue != "false" && unescapedQueryValue != "0" {
		panic("Unsupported value for Boolean query key: " + unescapedQueryValue)
	}

	return true, (unescapedQueryValue == "true" || unescapedQueryValue == "1")
}

func (this *BaseController) GetRequiredQueryValueBoolean(queryKey string) bool {
	hadVal, boolVal := this.GetOptionalQueryValueBoolean(queryKey, false)
	if !hadVal {
		panic("Required query key is empty/missing: " + queryKey)
	}

	return boolVal
}

func (this *BaseController) GetOptionalQueryValueString(queryKey, defaultVal string) (bool, string) {
	escapedQueryValue := this.Ctx.Input.Query(queryKey)
	if escapedQueryValue == "" {
		return false, defaultVal
	}
	unescapedQueryValue, err := url.QueryUnescape(escapedQueryValue)
	this.PanicIfError(err)
	return true, unescapedQueryValue
}

func (this *BaseController) GetRequiredQueryValueString(queryKey string) string {
	hadVal, unescapedQueryValue := this.GetOptionalQueryValueString(queryKey, "")
	if !hadVal {
		panic("Required query key is empty/missing: " + queryKey)
	}

	return unescapedQueryValue
}
