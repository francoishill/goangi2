package errorUtils

type ClientError struct {
	*BaseError
}

func (this *ClientError) SkipLogForError() bool {
	return true
}

func PanicClientError(errorFmt string, fmtParams ...interface{}) {
	errObj := &ClientError{
		BaseError: createBaseError(errorFmt, fmtParams...),
	}
	panic(errObj)
}
