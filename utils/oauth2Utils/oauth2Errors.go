package oauth2Utils

import (
	"github.com/RangelReale/osin"
)

const (
	E_INVALID_AUTH_DATA                          = "access_invalid_auth_data"
	E_ACCESS_EMPTY_USER_ERROR                    = "access_disallow_empty_user"
	E_ACCESS_DATA_MISSING_USER                   = "access_data_missing_user"
	E_INSUFFICIENT_SCOPE                         = "access_data_insufficient_scope"
	E_EMAIL_DOES_NOT_EXIST_OR_PASSWORD_INCORRECT = "email_not_exist_or_incorrect_password"
)

var errorMapKeys map[string]string = map[string]string{
	//Osin
	osin.E_INVALID_REQUEST:           "The request is missing a required parameter, includes an invalid parameter value, includes a parameter more than once, or is otherwise malformed.",
	osin.E_UNAUTHORIZED_CLIENT:       "The client is not authorized to request a token using this method.",
	osin.E_ACCESS_DENIED:             "The resource owner or authorization server denied the request.",
	osin.E_UNSUPPORTED_RESPONSE_TYPE: "The authorization server does not support obtaining a token using this method.",
	osin.E_INVALID_SCOPE:             "The requested scope is invalid, unknown, or malformed.",
	osin.E_SERVER_ERROR:              "The authorization server encountered an unexpected condition that prevented it from fulfilling the request.",
	osin.E_TEMPORARILY_UNAVAILABLE:   "The authorization server is currently unable to handle the request due to a temporary overloading or maintenance of the server.",
	osin.E_UNSUPPORTED_GRANT_TYPE:    "The authorization grant type is not supported by the authorization server.",
	osin.E_INVALID_GRANT:             "The provided authorization grant (e.g., authorization code, resource owner credentials) or refresh token is invalid, expired, revoked, does not match the redirection URI used in the authorization request, or was issued to another client.",
	osin.E_INVALID_CLIENT:            "Client authentication failed (e.g., unknown client, no client authentication included, or unsupported authentication method).",
	//Custom
	E_INVALID_AUTH_DATA:                          "Invalid access data, unable to authorize request",
	E_ACCESS_EMPTY_USER_ERROR:                    "Access tokens do not allow empty users, a valid user is required.",
	E_ACCESS_DATA_MISSING_USER:                   "Access data must have a valid user connected to it.",
	E_INSUFFICIENT_SCOPE:                         "This access token does not have sufficient rights in its scope.",
	E_EMAIL_DOES_NOT_EXIST_OR_PASSWORD_INCORRECT: "Email does not exist or the password is incorrect.",
}
