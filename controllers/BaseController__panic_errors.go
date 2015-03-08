package controllers

import (
	. "github.com/francoishill/goangi2/utils/errorUtils"
)

func (this *BaseController) PanicClientError(errorFmt string, fmtParams ...interface{}) {
	PanicClientError(errorFmt, fmtParams...)
}
