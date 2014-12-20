package oauth2Utils

func createOsinAuthorizeError(errorCode, errorString string) *OsinAuthorizeError {
	return &OsinAuthorizeError{
		ErrorCode:   errorCode,
		ErrorString: errorString,
	}
}

type OsinAuthorizeError struct {
	ErrorCode   string
	ErrorString string
}
