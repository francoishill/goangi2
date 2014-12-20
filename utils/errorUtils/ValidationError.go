package errorUtils

type ValidationError struct {
	*BaseError
}

func (this *ValidationError) SkipLogForError() bool {
	return true
}

func PanicValidationError(errorFmt string, fmtParams ...interface{}) {
	errObj := &ValidationError{
		BaseError: createBaseError(errorFmt, fmtParams...),
	}
	panic(errObj)
}
