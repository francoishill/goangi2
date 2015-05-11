package cookieUtils

import (
	"fmt"
	"strings"
)

const cREQUIRED_COOKIE_SECURITY_KEY_LENGTH = 16
const cREQUIRED_WEB_OAUTH_CLIENT_SECRET_LENGTH = 16

var DefaultCookieSecurityContext *CookieSecurityContext

type CookieSecurityContext struct {
	CookieExpireDays      int
	SecurityKey           string
	WebOauth2ClientId     string
	WebOauth2ClientSecret string
}

func CreateCookieSecurityContext(cookieExpireDays int, securityKey, webOauth2ClientId, webOauth2ClientSecret string) *CookieSecurityContext {
	if cookieExpireDays <= 0 {
		panic("Cookie expire days config must be > 0.")
	}
	if len(strings.Trim(securityKey, " ")) != cREQUIRED_COOKIE_SECURITY_KEY_LENGTH {
		panic(fmt.Sprintf("Security key length must be exactly %d characters long, key value '%s' is not.", cREQUIRED_COOKIE_SECURITY_KEY_LENGTH, securityKey))
	}
	if len(strings.Trim(webOauth2ClientSecret, " ")) != cREQUIRED_WEB_OAUTH_CLIENT_SECRET_LENGTH {
		panic(fmt.Sprintf("Web OAuth Client Secret length must be exactly %d characters long, secret value '%s' is not.", cREQUIRED_WEB_OAUTH_CLIENT_SECRET_LENGTH, webOauth2ClientSecret))
	}
	return &CookieSecurityContext{
		CookieExpireDays:      cookieExpireDays,
		SecurityKey:           securityKey,
		WebOauth2ClientId:     webOauth2ClientId,
		WebOauth2ClientSecret: webOauth2ClientSecret,
	}
}
