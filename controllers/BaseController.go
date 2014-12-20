package controllers

import (
	"github.com/astaxie/beego"

	. "github.com/francoishill/goangi2/context"
	. "github.com/francoishill/goangi2/utils/entityUtils"
	. "github.com/francoishill/goangi2/utils/errorUtils"
)

type BaseController struct {
	beego.Controller

	*BaseAppContext
}

func (this *BaseController) Prepare() {
	defer this.RecoverPanicAndServerError()
	this.Controller.Prepare()

	if DefaultBaseAppContext == nil {
		panic("Cannot use BaseController, DefaultBaseAppContext is nil")
	}

	this.BaseAppContext = DefaultBaseAppContext
}

func (this *BaseController) ServerJson_ErrorText(errorMessage string) {
	jsonData := map[string]string{
		"Error": errorMessage,
	}
	this.Data["json"] = jsonData
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
	this.ServerJson_ErrorText(userMessage)
}

func (this *BaseController) RecoverPanicAndServerError() {
	if r := recover(); r != nil {
		this.Ctx.Output.SetStatus(500)
		this.onAjaxRouterPanicRecovery(r)
	}
}

func (this *BaseController) CreateDefaultOrmContext() *OrmContext {
	return CreateOrmContext(this.Logger, nil, nil)
}
