package emailUtils

var HtmlEmailMessageType IEmailMessageType = &htmlEmailMessageType{}
var TextEmailMessageType IEmailMessageType = &textEmailMessageType{}

type EmailMessageTypeVisitor interface {
	visitTextEmailMessageType(visitor *textEmailMessageType)
	visitHtmlEmailMessageType(visitor *htmlEmailMessageType)
}

type IEmailMessageType interface {
	Accept(EmailMessageTypeVisitor)
}

type textEmailMessageType struct{}

func (this *textEmailMessageType) Accept(visitor EmailMessageTypeVisitor) {
	visitor.visitTextEmailMessageType(this)
}

type htmlEmailMessageType struct{}

func (this *htmlEmailMessageType) Accept(visitor EmailMessageTypeVisitor) {
	visitor.visitHtmlEmailMessageType(this)
}

type getContentTypeForEmailMessageTypeVisitor struct {
	ContentTypeString string
}

func (this *getContentTypeForEmailMessageTypeVisitor) visitTextEmailMessageType(visitor *textEmailMessageType) {
	this.ContentTypeString = "text/plain; charset=UTF-8"
}
func (this *getContentTypeForEmailMessageTypeVisitor) visitHtmlEmailMessageType(visitor *htmlEmailMessageType) {
	this.ContentTypeString = "text/html; charset=UTF-8"
}
