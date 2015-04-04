package validationUtils

import (
	. "github.com/francoishill/goangi2/utils/errorUtils"
)

func CheckBooleanCondition(condition bool, errorIfStringEmpty string) {
	if !condition {
		PanicValidationError(errorIfStringEmpty)
	}
}
