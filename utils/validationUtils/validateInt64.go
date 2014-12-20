package validationUtils

import (
	"fmt"

	. "github.com/francoishill/goangi2/utils/errorUtils"
)

func CheckInt64InRange(intToCheck, min, max int64, errorIfOutOfRange string) {
	if intToCheck < min || intToCheck > max {
		PanicValidationError(errorIfOutOfRange + fmt.Sprintf(" (range is %d-%d)", min, max))
	}
}

func CheckInt64IsAtLeast(intToCheck, min int64, errorIfNotIsAtLeast string) {
	if intToCheck < min {
		PanicValidationError(errorIfNotIsAtLeast + fmt.Sprintf(" (minimum %d)", min))
	}
}
