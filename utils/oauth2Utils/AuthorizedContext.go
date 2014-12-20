package oauth2Utils

import (
	. "github.com/francoishill/goangi2/utils/encodingUtils"
)

type tmpSetSecureCookieInterface interface {
	SetSecureCookie(secret, name, value string, others ...interface{})
}

type tmpGetSecureCookieInterface interface {
	GetSecureCookie(secret, key string) (string, bool)
}

func CreateAuthorizedContext(user IExpectedUser, scope string, accessToken string) *AuthorizedContext {
	return &AuthorizedContext{
		User:        user,
		Scope:       scope,
		accessToken: accessToken,
	}
}

type AuthorizedContext struct {
	User        IExpectedUser
	Scope       string
	accessToken string
}

func (this *AuthorizedContext) GetAccessToken() string {
	return this.accessToken
}

func (this *AuthorizedContext) getCookieSecret() string {
	secret := EncodeMd5(this.User.GetRands() + this.User.GetPassword())
	return secret
}

func TempCookieSecret(randsPlusPassword string) string {
	secret := EncodeMd5(randsPlusPassword)
	return secret
}

func (this *AuthorizedContext) SetCookie(router tmpSetSecureCookieInterface, name, value string) {
	totalDays := 7
	totalSeconds := int64(86400 * totalDays)
	router.SetSecureCookie(this.getCookieSecret(), name, value, totalSeconds)
}

func (this *AuthorizedContext) GetCookie(router tmpGetSecureCookieInterface, name string) (string, bool) {
	return router.GetSecureCookie(this.getCookieSecret(), name)
}
