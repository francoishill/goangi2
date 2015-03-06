package cookieUtils

import (
	"fmt"
	"github.com/astaxie/beego/context"
	"strconv"
	"strings"
	"time"
)

func scrambleSaltString(normalSaltString string) string {
	return normalSaltString[:12] + DefaultCookieSecurityContext.SecurityKey + normalSaltString[10:] //Yes some characters are repeated
}

func setEncryptedCookie(ctx *context.Context, key, value string) {
	expireDays := 7
	totalExpireSeconds := int64(86400 * expireDays)

	userAgent := ctx.Request.UserAgent()

	//Set the 'dat1' cookie as the date string in the original format + '_length_of_user_agent'
	normalSaltString := time.Now().Format("2006_01_02__15_04_05") + fmt.Sprintf("_%d", len(userAgent))
	ctx.SetCookie(cSALT_SALT_COOKIE_KEY, normalSaltString, totalExpireSeconds)

	valueWithUserAgentPrepended := userAgent + value
	//Now encrypt the cookie with a secret containing our COOKIE_SECRET as well as a scrambled date string
	secretValueScrambled := scrambleSaltString(normalSaltString)
	ctx.SetSecureCookie(secretValueScrambled, key, valueWithUserAgentPrepended, totalExpireSeconds)
}

func SetEncryptedAccessTokenInCookie(ctx *context.Context, accessToken string) {
	setEncryptedCookie(ctx, cUSER_AGENT_AND_ACCESS_TOKEN_COOKIE_KEY_NAME, accessToken)
}

func DeleteAccessTokenCookies(ctx *context.Context) {
	ctx.SetCookie(cSALT_SALT_COOKIE_KEY, "")
	ctx.SetCookie(cUSER_AGENT_AND_ACCESS_TOKEN_COOKIE_KEY_NAME, "")
}

func getDecryptedCookie(ctx *context.Context, key string) (string, bool) {
	normalSaltString := ctx.GetCookie(cSALT_SALT_COOKIE_KEY)
	if strings.Trim(normalSaltString, " ") == "" {
		return "", false
	}
	lengthOfUserAgent, err := strconv.ParseInt(normalSaltString[strings.LastIndex(normalSaltString, "_")+1:], 10, 32)
	if err != nil {
		panic(err)
	}

	secretValueScrambled := scrambleSaltString(normalSaltString)
	decryptedValueWithUserAgentPrefix, success := ctx.GetSecureCookie(secretValueScrambled, key)
	if !success {
		return "", false
	}
	return decryptedValueWithUserAgentPrefix[lengthOfUserAgent:], true
}

func GetDecryptedAccessTokenFromCookie(ctx *context.Context) (string, bool) {
	return getDecryptedCookie(ctx, cUSER_AGENT_AND_ACCESS_TOKEN_COOKIE_KEY_NAME)
}
