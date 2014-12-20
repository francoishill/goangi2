package validationUtils

import (
	"strings"
)

func CheckStringNotEmpty(strToCheck, errorIfStringEmpty string) {
	if strings.Trim(strToCheck, " ") == "" {
		panic(errorIfStringEmpty)
	}
}
