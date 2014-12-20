package controllers

import (
	"net/url"
	"strconv"
	"strings"
)

func (this *BaseController) GetOptionalParamBoolean(paramName string, defaultVal bool) (hadParam bool, boolVal bool) {
	escapedParam := this.Ctx.Input.Param(paramName)
	if escapedParam == "" {
		return false, defaultVal
	}
	unescapedParam, err := url.QueryUnescape(escapedParam)
	this.PanicIfError(err)
	unescapedParam = strings.ToLower(strings.Trim(unescapedParam, " "))

	if unescapedParam != "true" && unescapedParam != "1" &&
		unescapedParam != "false" && unescapedParam != "0" {
		panic("Unsupported value for Boolean query parameter: " + unescapedParam)
	}

	return true, (unescapedParam == "true" || unescapedParam == "1")
}

func (this *BaseController) GetRequiredParamBoolean(paramName string) bool {
	hadVal, boolVal := this.GetOptionalParamBoolean(paramName, false)
	if !hadVal {
		panic("Required query parameter is empty/missing: " + paramName)
	}

	return boolVal
}

func (this *BaseController) GetOptionalParamString(paramName, defaultVal string) (bool, string) {
	escapedParam := this.Ctx.Input.Param(paramName)
	if escapedParam == "" {
		return false, defaultVal
	}
	unescapedParam, err := url.QueryUnescape(escapedParam)
	this.PanicIfError(err)
	return true, unescapedParam
}

func (this *BaseController) GetRequiredParamString(paramName string) string {
	hadVal, unescapedParam := this.GetOptionalParamString(paramName, "")
	if !hadVal {
		panic("Required query parameter is empty/missing: " + paramName)
	}

	return unescapedParam
}

func (this *BaseController) GetOptionalParamInt64(paramName string, defaultVal int64) (bool, int64) {
	escapedParam := this.Ctx.Input.Param(paramName)
	if escapedParam == "" {
		return false, defaultVal
	}
	unescapedParam, err := url.QueryUnescape(escapedParam)
	this.PanicIfError(err)

	intVal, err := strconv.ParseInt(unescapedParam, 10, 64)
	this.PanicIfError(err)

	return true, intVal
}

func (this *BaseController) GetRequiredParamInt64(paramName string) int64 {
	hadVal, int64Val := this.GetOptionalParamInt64(paramName, -1)
	if !hadVal {
		panic("Required query parameter is empty/missing: " + paramName)
	}

	return int64Val
}
