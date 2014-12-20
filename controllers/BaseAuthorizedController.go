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
	defer this.RecoverPanicAndServerError_InControllerPrepare()

	this.BaseController.Prepare()

	this.OsinResponse = OsinServerObject.NewResponse()

	this.AuthorizedContext = GetAuthorizedContextFromAccessToken(this.OsinResponse, this.Ctx)
	if this.AuthorizedContext == nil {
		panic("Internal server error [code 600001]")
	}

	//TODO: Not handling cookies at this point
}

func (this *BaseAuthorizedController) RecoverPanicAndServerError() {
	defer this.BaseController.RecoverPanicAndServerError() //To catch the non-osin errors

	if r := recover(); r != nil {
		switch e := r.(type) {
		case *OsinAuthorizeError:
			if strings.EqualFold(e.ErrorCode, E_INVALID_AUTH_DATA) {
				DeleteAccessTokenCookies(this.Controller.Ctx)
			}
			this.OsinResponse.ErrorStatusCode = 401
			this.OsinResponse.SetError(e.ErrorCode, e.ErrorString)
			OverwriteOsinResponseErrorWithOwn(this.OsinResponse)
			osin.OutputJSON(this.OsinResponse, this.Ctx.ResponseWriter, this.Ctx.Request)
		default:
			panic(r) //So it can be caught by the Base Controller
		}
	}
}

func (this *BaseAuthorizedController) CreateDefaultOrmContext() *OrmContext {
	return CreateOrmContext(this.Logger, nil, this.AuthorizedContext.User)
}
