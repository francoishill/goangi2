package jsonUtils

import (
	"encoding/json"
	. "github.com/francoishill/goangi2/utils/errorUtils"
)

type iValidate interface {
	Validate()
}

func UnmarshalFromJsonAndValidate(bodyBytes []byte, outputObject iValidate) {
	err := json.Unmarshal(bodyBytes, outputObject)
	if err != nil {
		PanicRequestJsonUnmarshalError(bodyBytes, outputObject, "Unable to marshal bodyBytes (%s) to object (%#v, %T), ERROR: %+v",
			string(bodyBytes), outputObject, outputObject, err)
	}
	outputObject.Validate()
}
