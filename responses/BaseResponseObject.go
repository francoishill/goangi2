package responses

type IRouterResponseObject interface {
	GetSuccess() bool
	GetMessageKey() string
}

type BaseResponseObject struct {
	Success    bool
	MessageKey string
}

func CreateBaseResponseObject(success bool, messageKey string) *BaseResponseObject {
	return &BaseResponseObject{
		Success:    success,
		MessageKey: messageKey,
	}
}

func (this *BaseResponseObject) GetSuccess() bool {
	return this.Success
}

func (this *BaseResponseObject) SetSuccess(newSuccess bool) {
	this.Success = newSuccess
}

func (this *BaseResponseObject) GetMessageKey() string {
	return this.MessageKey
}

func (this *BaseResponseObject) SetMessageKey(newMessageKey string) {
	this.MessageKey = newMessageKey
}
