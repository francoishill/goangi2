package errors

import (
	"fmt"
)

type NotLoggedError struct {
	ErrorString string
}

func PanicNotLoggedError(errorFmt string, fmtParams ...interface{}) {
	var errorStr string
	if len(fmtParams) > 0 {
		errorStr = fmt.Sprintf(errorFmt, fmtParams...)
	} else {
		errorStr = errorFmt
	}

	errObj := &NotLoggedError{
		ErrorString: errorStr,
	}
	panic(errObj)
}
