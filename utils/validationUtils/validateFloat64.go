package validationUtils

import (
	"fmt"

	. "github.com/francoishill/goangi2/utils/errorUtils"
)

func CheckFloat64InRange(intToCheck, min, max float64, errorIfOutOfRange string) {
	if intToCheck < min || intToCheck > max {
		PanicValidationError(errorIfOutOfRange + fmt.Sprintf(" (range is %d-%d)", min, max))
	}
}

func CheckFloat64Is_MinimumInclusive(intToCheck, minInclusive float64, errorIfNotIsAtLeast string) {
	if intToCheck < minInclusive {
		PanicValidationError(errorIfNotIsAtLeast + fmt.Sprintf(" (minimumInclusive %d)", minInclusive))
	}
}

func CheckFloat64Is_MinimumExclusive(intToCheck, min float64, errorIfNotIsAtLeast string) {
	if intToCheck <= min {
		PanicValidationError(errorIfNotIsAtLeast + fmt.Sprintf(" (minimumExclusive %d)", min))
	}
}
