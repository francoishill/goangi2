package errorUtils

type RequestJsonUnmarshalError struct {
	*BaseError
	OutputObject      interface{} //The object we are trying to unmarshal to
	BodyBytesAsString string
}

func (this *RequestJsonUnmarshalError) SkipLogForError() bool {
	return true
}

func PanicRequestJsonUnmarshalError(bodyBytes []byte, outputObject interface{}, errorFmt string, fmtParams ...interface{}) {
	errObj := &RequestJsonUnmarshalError{
		BaseError:         createBaseError(errorFmt, fmtParams...),
		OutputObject:      outputObject,
		BodyBytesAsString: string(bodyBytes),
	}
	panic(errObj)
}
