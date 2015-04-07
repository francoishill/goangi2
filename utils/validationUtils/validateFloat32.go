package validationUtils

import (
	"fmt"

	. "github.com/francoishill/goangi2/utils/errorUtils"
)

func CheckFloat32InRange(intToCheck, min, max float32, errorIfOutOfRange string) {
	if intToCheck < min || intToCheck > max {
		PanicValidationError(errorIfOutOfRange + fmt.Sprintf(" (range is %d-%d)", min, max))
	}
}

func CheckFloat32Is_MinimumInclusive(intToCheck, minInclusive float32, errorIfNotIsAtLeast string) {
	if intToCheck < minInclusive {
		PanicValidationError(errorIfNotIsAtLeast + fmt.Sprintf(" (minimumInclusive %d)", minInclusive))
	}
}

func CheckFloat32Is_MinimumExclusive(intToCheck, min float32, errorIfNotIsAtLeast string) {
	if intToCheck <= min {
		PanicValidationError(errorIfNotIsAtLeast + fmt.Sprintf(" (minimumExclusive %d)", min))
	}
}
