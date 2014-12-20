package validationUtils

import (
	"regexp"
	"strings"

	. "github.com/francoishill/goangi2/utils/errorUtils"
)

//Thanks to github.com/astaxie/beego/validation
var emailPattern = regexp.MustCompile("[\\w!#$%&'*+/=?^_`{|}~-]+(?:\\.[\\w!#$%&'*+/=?^_`{|}~-]+)*@(?:[\\w](?:[\\w-]*[\\w])?\\.)+[a-zA-Z0-9](?:[\\w-]*[\\w])?")

func IsValidEmail(emailStr string) bool {
	return emailPattern.Match([]byte(emailStr))
}

func CheckValidEmail(emailStr, errorIfInvalidEmail string) {
	if strings.Trim(emailStr, " ") == "" {
		PanicValidationError(errorIfInvalidEmail)
	}
	if !IsValidEmail(emailStr) {
		PanicValidationError(errorIfInvalidEmail)
	}
}
