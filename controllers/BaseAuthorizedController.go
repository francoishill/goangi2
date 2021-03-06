package controllers

import (
	"github.com/RangelReale/osin"
	"strings"

	. "github.com/francoishill/goangi2/utils/cookieUtils"
	. "github.com/francoishill/goangi2/utils/entityUtils"
	. "github.com/francoishill/goangi2/utils/oauth2Utils"
)

type BaseAuthorizedController struct {
	BaseController

	OsinResponse *osin.Response
	*AuthorizedContext
}

func (this *BaseAuthorizedController) Prepare() {
	defer this.RecoverPanicAndServeError_InControllerPrepare()
	this.BaseController.Prepare()

	this.OsinResponse = OsinServerObject.NewResponse()

	accessTokenFromCookie, successGotTokenFromCookie := GetDecryptedAccessTokenFromCookie(this.Ctx)
	if successGotTokenFromCookie {
		this.Ctx.Request.ParseForm()
		this.Ctx.Request.Form.Set("code", accessTokenFromCookie)
	}

	this.AuthorizedContext = GetAuthorizedContextFromAccessToken(this.OsinResponse, this.Ctx)
	if this.AuthorizedContext == nil {
		panic("Internal server error [code 600001]")
	}

	//TODO: Not handling cookies at this point
}

func (this *BaseAuthorizedController) Finish() {
	if this.OsinResponse != nil {
		this.OsinResponse.Close()
	}
}

func (this *BaseAuthorizedController) BaseAuthorizedController_PanicInvalidAuthData() {
	PanicInvalidAuthData()
}

func (this *BaseAuthorizedController) RecoverPanicAndServeError() {
	defer this.BaseController.RecoverPanicAndServeError() //BaseController to catch the non-osin errors

	if r := recover(); r != nil {
		switch e := r.(type) {
		case *OsinAuthorizeError:
			//TODO: I have a suspicion this block (OsinAuthorizeError) should never be reached. We handle OAuth stuff in the Prepare, so therefore gets caught in RecoverPanicAndServeError_InControllerPrepare

			if strings.EqualFold(e.ErrorCode, E_INVALID_AUTH_DATA) {
				DeleteAccessTokenCookies(this.Controller.Ctx)
				ClearUserCookies(this.Controller.Ctx)
			}
			this.OsinResponse.ErrorStatusCode = 401
			this.OsinResponse.SetError(e.ErrorCode, e.ErrorString)
			OverwriteOsinResponseErrorWithOwn(this.OsinResponse)
			break
		}
		panic(r) //So it can be caught by the Base Controller
	}
}

func (this *BaseAuthorizedController) RecoverPanicAndServeError_InControllerPrepare() {
	defer this.BaseController.RecoverPanicAndServeError_InControllerPrepare() //BaseController to catch the non-osin errors

	if r := recover(); r != nil {
		//Serve the error as-is, otherwise the osin errors will
		switch e := r.(type) {
		case *OsinAuthorizeError:
			//This does not work correctly if the GZip is on
			this.Controller.Ctx.Output.EnableGzip = false
			if strings.EqualFold(e.ErrorCode, E_INVALID_AUTH_DATA) {
				DeleteAccessTokenCookies(this.Controller.Ctx)
				ClearUserCookies(this.Controller.Ctx)
			}
			this.OsinResponse.ErrorStatusCode = 401
			this.OsinResponse.SetError(e.ErrorCode, e.ErrorString)
			OverwriteOsinResponseErrorWithOwn(this.OsinResponse)
			break
		}
		panic(r)
	}
}

func (this *BaseAuthorizedController) CreateDefaultRouterOrmContext(beginTransaction bool) *OrmContext {
	return CreateOrmContext(this.Logger, nil, beginTransaction)
}
