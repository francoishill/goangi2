package cookieUtils

//Use salt to encrypt the ats

const (
	cUSER_ID_COOKIE_KEY_NAME = "uid" //Keep in sync with clientside code

	cSALT_SALT_COOKIE_KEY                        = "dat1" //Salt used, is the date and length of user agent
	cUSER_AGENT_AND_ACCESS_TOKEN_COOKIE_KEY_NAME = "dat2" //Encrypted version of useragent+accessToken
)
