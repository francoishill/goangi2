package cookieUtils

import (
	"fmt"
	"github.com/astaxie/beego/context"
)

func SetUserCookies(ctx *context.Context, uid int64) {
	expireDays := 14
	totalExpireSeconds := int64(86400 * expireDays)
	ctx.SetCookie(cUSER_ID_COOKIE_KEY_NAME, fmt.Sprintf("%d", uid), totalExpireSeconds)
}

func ClearUserCookies(ctx *context.Context) {
	ctx.SetCookie(cUSER_ID_COOKIE_KEY_NAME, "")
}
