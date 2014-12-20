package validationUtils

import (
	"strings"

	. "github.com/francoishill/goangi2/utils/errorUtils"
)

func CheckStringNotEmpty(strToCheck, errorIfStringEmpty string) {
	if strings.Trim(strToCheck, " ") == "" {
		PanicValidationError(errorIfStringEmpty)
	}
}
