package validationUtils

import (
	"regexp"

	. "github.com/francoishill/goangi2/utils/errorUtils"
)

//Thanks to github.com/astaxie/beego/validation
var emailPattern = regexp.MustCompile("[\\w!#$%&'*+/=?^_`{|}~-]+(?:\\.[\\w!#$%&'*+/=?^_`{|}~-]+)*@(?:[\\w](?:[\\w-]*[\\w])?\\.)+[a-zA-Z0-9](?:[\\w-]*[\\w])?")

func IsValidEmail(emailStr string) bool {
	return emailPattern.Match([]byte(emailStr))
}

func CheckValidEmail(emailStr, errorIfInvalidEmail string) {
	if !IsValidEmail(emailStr) {
		PanicValidationError(errorIfInvalidEmail)
	}
}
