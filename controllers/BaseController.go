package controllers

import (
	"fmt"
	"github.com/astaxie/beego"

	. "github.com/francoishill/goangi2/context"
	. "github.com/francoishill/goangi2/responses"
	. "github.com/francoishill/goangi2/utils/entityUtils"
	. "github.com/francoishill/goangi2/utils/errorUtils"
	. "github.com/francoishill/goangi2/utils/oauth2Utils"
)

type BaseController struct {
	beego.Controller

	*BaseAppContext
}

func (this *BaseController) Prepare() {
	defer this.RecoverPanicAndServeError_InControllerPrepare()

	this.Controller.Prepare()

	if DefaultBaseAppContext == nil {
		panic("Cannot use BaseController, DefaultBaseAppContext is nil")
	}

	this.BaseAppContext = DefaultBaseAppContext
}

func (this *BaseController) PanicIfError(err error) {
	if err != nil {
		panic(err)
	}
}

func (this *BaseController) ServeJson_ErrorText(errorMessage string) {
	jsonData := map[string]interface{}{
		"Success": false,
		"Error":   errorMessage,
	}
	this.Data["json"] = jsonData
	this.ServeJson()
}

func (this *BaseController) ServeJson_SuccessText(successMessage string) {
	jsonData := map[string]interface{}{
		"Success": true,
	}
	if successMessage != "" {
		jsonData["Message"] = successMessage
	}
	this.Data["json"] = jsonData
	this.ServeJson()
}

func (this *BaseController) ServeJsonResponseObject(responseObject IRouterResponseObject) {
	this.Data["json"] = responseObject
	this.ServeJson()
}

func (this *BaseController) onAjaxRouterPanicRecovery(recoveryObj interface{}) {
	this.Ctx.Output.SetStatus(500)
	requestUrl := this.Ctx.Request.URL
	remoteAddress := this.Ctx.Request.RemoteAddr
	userAgent := this.Ctx.Input.UserAgent()

	userMessage := LogRouterError_And_ExtractUserMessage(
		this.BaseAppContext.Logger,
		"Controller error:",
		requestUrl,
		remoteAddress,
		this.Ctx.Input.Proxy(),
		userAgent,
		recoveryObj,
	)
	this.ServeJson_ErrorText(userMessage)
}

func (this *BaseController) RecoverPanicAndServeError() {
	if r := recover(); r != nil {
		this.Ctx.Output.SetStatus(500)
		this.onAjaxRouterPanicRecovery(r)
	}
}

func (this *BaseController) RecoverPanicAndServeError_InControllerPrepare() {
	if r := recover(); r != nil {
		//This does not work correctly if the GZip is on
		this.Controller.Ctx.Output.EnableGzip = false

		this.Ctx.Output.SetStatus(500)
		//Serve the error as-is, otherwise the osin errors will
		switch e := r.(type) {
		case *OsinAuthorizeError:
			this.Data["json"] = e
			this.ServeJson()
		case string:
			this.ServeJson_ErrorText(e)
		case error:
			this.ServeJson_ErrorText(e.Error())
		default:
			this.ServeJson_ErrorText(fmt.Sprintf("%+v", r))
		}
		this.StopRun()
	}
}

func (this *BaseController) CreateDefaultRouterOrmContext(beginTransaction bool) *OrmContext {
	return CreateOrmContext(this.Logger, nil, beginTransaction)
}
