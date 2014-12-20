package errorUtils

import (
	"fmt"
)

type BaseError struct {
	ErrorString string
}

func (this *BaseError) PrettyErrorSPrint() string {
	return this.ErrorString
}

func createBaseError(errorFmt string, fmtParams ...interface{}) *BaseError {
	var errorStr string
	if len(fmtParams) > 0 {
		errorStr = fmt.Sprintf(errorFmt, fmtParams...)
	} else {
		errorStr = errorFmt
	}

	return &BaseError{
		ErrorString: errorStr,
	}
}
