package oauth2Utils

import (
	"github.com/RangelReale/osin"
	"github.com/astaxie/beego/context"
	"net/http"
	"strings"

	. "github.com/francoishill/goangi2/utils/cookieUtils"
)

const (
	cOSIN_ACCESS_OUTPUT_ACCESS_TOKEN_MAP_KEY = "access_token"
)

type StringPredicate func(string) bool

var OsinServerObject *osin.Server

type IExpectedUser interface {
	GetRands() string
	GetPassword() string
	GetId() int64
	IAmAUser()
}

type iAuthUserProvider interface {
	DoVerifyUser(userName, password string) (bool, IExpectedUser) //This can handle both login/register
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

func InjectCodeIntoFormIfWasPassedViaAuthorizationHeader(request *http.Request) {
	authorizationHeader := request.Header.Get("Authorization")
	accessTokenFromAuthHeader := ""
	if strings.HasPrefix(authorizationHeader, "Bearer ") {
		//This means if the "code" was set in the get/set form it will be overwritten if we use the Authorization header
		accessTokenFromAuthHeader = authorizationHeader[7:]
		request.ParseForm()
		request.Form.Set("code", accessTokenFromAuthHeader)
	}
}

func OverwriteOsinResponseErrorWithOwn(osinResponse *osin.Response) {
	errKey, ok := osinResponse.Output["error"].(string)
	if ok {
		if errKey == osin.E_INVALID_REQUEST {
			osinResponse.SetError(errKey, errorMapKeys[errKey])
		}
	}
}

func OverwriteOsinResponseErrorWithOwn_SpecifyErrorKey(osinResponse *osin.Response, errorKey string) {
	osinResponse.SetError(errorKey, errorMapKeys[errorKey])
}

func GetAuthorizedContextFromAccessToken(osinResponse *osin.Response, ctx *context.Context) *AuthorizedContext {
	InjectCodeIntoFormIfWasPassedViaAuthorizationHeader(ctx.Request)

	var usr IExpectedUser
	var castedToUserType bool

	ir := OsinServerObject.HandleInfoRequest(osinResponse, ctx.Request)
	if ir == nil {
		panic(createOsinAuthorizeError(E_INVALID_AUTH_DATA, errorMapKeys[E_INVALID_AUTH_DATA]))
	}

	if ir.AccessData.UserData == nil {
		panic(createOsinAuthorizeError(E_ACCESS_DATA_MISSING_USER, errorMapKeys[E_ACCESS_DATA_MISSING_USER]+"(1)"))
	}
	if strings.Trim(ir.AccessData.AccessToken, " ") == "" {
		panic(createOsinAuthorizeError(E_ACCESS_DATA_MISSING_USER, errorMapKeys[E_ACCESS_DATA_MISSING_USER]+"(2)"))
	}
	usr, castedToUserType = ir.AccessData.UserData.(IExpectedUser)
	if !castedToUserType {
		panic(createOsinAuthorizeError(E_ACCESS_DATA_MISSING_USER, errorMapKeys[E_ACCESS_DATA_MISSING_USER]+"(3)"))
	}

	return CreateAuthorizedContext(usr, ir.AccessData.Scope, ir.AccessData.AccessToken)
}

func CheckRequiredScopeSatisfied(responseWriter http.ResponseWriter, authorizedScope string, functionToCheckRequiredScope StringPredicate) {
	if !functionToCheckRequiredScope(authorizedScope) {
		panic(createOsinAuthorizeError(E_INSUFFICIENT_SCOPE, errorMapKeys[E_INSUFFICIENT_SCOPE]))
	}
}

func ServeAccessTokenWithRouter(ctx *context.Context) {
	w := ctx.ResponseWriter
	r := ctx.Request

	InjectCodeIntoFormIfWasPassedViaAuthorizationHeader(r)

	resp := OsinServerObject.NewResponse()
	if ir := OsinServerObject.HandleInfoRequest(resp, r); ir != nil {
		OsinServerObject.FinishInfoRequest(resp, r, ir)
	}

	if resp.IsError {
		OverwriteOsinResponseErrorWithOwn(resp)
	}
	osin.OutputJSON(resp, w, r)
}

func setExpirationForAccessRequest(accessRequest *osin.AccessRequest) {
	accessRequest.Expiration = 60 * 60 * 24 * 365 * 10 //Ten years
	//accessRequest.Expiration = 60 * 60 * 24 * 365 * 1 //One year
	// accessRequest.Expiration = 60 * 60 * 24 //One day
	//beego.Warning("Currently hardcoded the access token expiration to one year, is this correct?")
}

func ExtractAccessTokenFromSuccessfulResponseData(responseData osin.ResponseData) (string, bool) {
	if token, ok := responseData[cOSIN_ACCESS_OUTPUT_ACCESS_TOKEN_MAP_KEY]; ok {
		if stringToken, ok := token.(string); ok {
			return stringToken, true
		}
	}

	return "", false
}

type outputHandlerFunc func(ctx *context.Context) bool

func AuthorizeAndServeNewAccessTokenWithRouter(ctx *context.Context, authUserProvider iAuthUserProvider, setCookies bool, successfulOutputHandler outputHandlerFunc) {
	resp := OsinServerObject.NewResponse()
	r := ctx.Request
	w := ctx.ResponseWriter

	ar := OsinServerObject.HandleAccessRequest(resp, r)
	if ar != nil {
		switch ar.Type {
		/*case osin.AUTHORIZATION_CODE:
		ar.Authorized = true*/
		case osin.REFRESH_TOKEN:
			ar.Authorized = true
			setExpirationForAccessRequest(ar)
		case osin.PASSWORD:
			ar.Authorized, ar.UserData = authUserProvider.DoVerifyUser(ar.Username, ar.Password)
			if !ar.Authorized {
				OverwriteOsinResponseErrorWithOwn_SpecifyErrorKey(resp, E_EMAIL_DOES_NOT_EXIST_OR_PASSWORD_INCORRECT)
			}

			/*case osin.CLIENT_CREDENTIALS:
			ar.Authorized = true*/
		}
		OsinServerObject.FinishAccessRequest(resp, r, ar)
	}

	if resp.IsError {
		if resp.InternalError != nil {
			OverwriteOsinResponseErrorWithOwn(resp)
		}

		resp.ErrorStatusCode = 401
		resp.StatusCode = 401
	}

	if setCookies && !resp.IsError {
		if accessToken, ok := ExtractAccessTokenFromSuccessfulResponseData(resp.Output); ok {
			if ar != nil && ar.UserData != nil {
				if usr, ok := ar.UserData.(IExpectedUser); ok {
					SetUserCookies(ctx, usr.GetId())
				}
			}

			SetEncryptedAccessTokenInCookie(ctx, accessToken)
		}
	}

	if successfulOutputHandler != nil {
		if handledOutput := successfulOutputHandler(ctx); handledOutput {
			return
		}
	}

	err := osin.OutputJSON(resp, w, r)
	checkError(err)
}

func InitOsinServerObject() {
	sconfig := osin.NewServerConfig()

	sconfig.AllowGetAccessRequest = false
	sconfig.AllowedAuthorizeTypes = osin.AllowedAuthorizeType{osin.CODE, osin.TOKEN}
	sconfig.AllowedAccessTypes = osin.AllowedAccessType{
		// osin.AUTHORIZATION_CODE,
		osin.REFRESH_TOKEN,
		osin.PASSWORD,
		// osin.CLIENT_CREDENTIALS,
	}
	OsinServerObject = osin.NewServer(sconfig, NewOAuth2Storage())
}
