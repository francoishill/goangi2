package controllers

import (
	"encoding/json"
	. "github.com/francoishill/goangi2/requests"
)

func (this *BaseController) UnmarshalJsonRequestObjectFromBody(requestObject IRequestObject) {
	err := json.Unmarshal(this.Ctx.Input.RequestBody, requestObject)
	if err != nil {
		panic(err)
	}
	requestObject.Validate()
}
