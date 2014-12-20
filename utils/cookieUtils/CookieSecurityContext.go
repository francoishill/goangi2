package cookieUtils

import (
	"fmt"
	"strings"
)

const cREQUIRED_COOKIE_SECURITY_KEY_LENGTH = 16

var DefaultCookieSecurityContext *CookieSecurityContext

type CookieSecurityContext struct {
	SecurityKey string
}

func CreateCookieSecurityContext(securityKey string) *CookieSecurityContext {
	if len(strings.Trim(securityKey, " ")) != cREQUIRED_COOKIE_SECURITY_KEY_LENGTH {
		panic(fmt.Sprintf("Security key length must be exactly %d characters long, key value '%s' is not.", cREQUIRED_COOKIE_SECURITY_KEY_LENGTH, securityKey))
	}
	return &CookieSecurityContext{
		SecurityKey: securityKey,
	}
}
