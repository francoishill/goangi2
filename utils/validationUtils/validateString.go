package validationUtils

import (
	"fmt"
	"strings"

	. "github.com/francoishill/goangi2/utils/errorUtils"
)

func CheckStringNotEmpty(strToCheck, errorIfStringEmpty string) {
	if strings.Trim(strToCheck, " ") == "" {
		PanicValidationError(errorIfStringEmpty)
	}
}

func CheckStringIsOneOf(strIn string, oneOf []string, ignoreCase, doTrim bool, errorIfNotInList string) {
	for _, o := range oneOf {
		str := strIn
		one := o

		if doTrim {
			str = strings.Trim(strIn, " \t")
			one = strings.Trim(o, " \t")
		}

		var equal bool
		if ignoreCase {
			equal = strings.EqualFold(str, one)
		} else {
			equal = str == one
		}

		if equal {
			return
		}
	}

	PanicValidationError(errorIfNotInList + fmt.Sprintf(" (must be in list %+v)", oneOf))
}
