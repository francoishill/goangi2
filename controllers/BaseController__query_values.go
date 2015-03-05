package controllers

import (
	"net/url"
	"strconv"
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

func (this *BaseController) GetOptionalQueryInt64(queryName string, defaultVal int64) (bool, int64) {
	escapedQuery := this.Ctx.Input.Query(queryName)
	if escapedQuery == "" {
		return false, defaultVal
	}
	unescapedQuery, err := url.QueryUnescape(escapedQuery)
	this.PanicIfError(err)

	intVal, err := strconv.ParseInt(unescapedQuery, 10, 64)
	this.PanicIfError(err)

	return true, intVal
}

func (this *BaseController) GetRequiredQueryInt64(queryName string) int64 {
	hadVal, int64Val := this.GetOptionalQueryInt64(queryName, -1)
	if !hadVal {
		panic("Required query is empty/missing: " + queryName)
	}

	return int64Val
}

func (this *BaseController) GetOptionalQueryInt64CsvArray(queryName string, defaultVal []int64) (bool, []int64) {
	escapedQuery := this.Ctx.Input.Query(queryName)
	if escapedQuery == "" {
		return false, defaultVal
	}
	unescapedQuery, err := url.QueryUnescape(escapedQuery)
	this.PanicIfError(err)

	intSlice := []int64{}

	intStrings := strings.Split(unescapedQuery, ",")
	for _, intStr := range intStrings {
		intVal, err := strconv.ParseInt(intStr, 10, 64)
		this.PanicIfError(err)
		intSlice = append(intSlice, intVal)
	}

	return true, intSlice
}

func (this *BaseController) GetRequiredQueryInt64CsvArray(queryName string) []int64 {
	hadVal, int64Slice := this.GetOptionalQueryInt64CsvArray(queryName, []int64{})
	if !hadVal {
		panic("Required query is empty/missing: " + queryName)
	}

	return int64Slice
}
