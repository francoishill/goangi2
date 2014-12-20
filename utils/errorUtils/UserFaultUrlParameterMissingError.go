package errorUtils

type UserFaultUrlParameterMissingError struct {
	*BaseError
}

func (this *UserFaultUrlParameterMissingError) SkipLogForError() bool {
	return true
}

func PanicUserFaultUrlParameterMissingError(errorFmt string, fmtParams ...interface{}) {
	errObj := &UserFaultUrlParameterMissingError{
		BaseError: createBaseError(errorFmt, fmtParams...),
	}
	panic(errObj)
}
