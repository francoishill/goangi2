package emailUtils

import (
	"time"
)

type iEmailView interface {
	GetBaseViewData() interface{}
	GetBaseViewDataName() string
	GetViewDataName() string
	GetEmailTemplatePath() string
}

func EnqueueEmail(
	sendGridAdvancedSuppressionManagerGroup int, //Need to find a more generic way to handle it so we do not need to pass it in for non-sendgrid email providers
	viewData iEmailView,
	emailContext *EmailContext,
	scheduleSendMailNow bool,
	sendDueTime time.Time,
	to []*EmailRecipient,
	from *EmailRecipient,
	subject string,
	msgType IEmailMessageType,
	debugInfo string,
	attachments ...*EmailAttachment) {

	if viewData.GetViewDataName() == "" || viewData.GetBaseViewDataName() == "" {
		panic("Both 'ViewDataName' and 'BaseViewDataName' must be non-empty")
	}

	templatePath := viewData.GetEmailTemplatePath()
	emailDataObject := make(map[interface{}]interface{})
	emailDataObject[viewData.GetBaseViewDataName()] = viewData.GetBaseViewData()
	emailDataObject[viewData.GetViewDataName()] = viewData
	emailBody := RenderGoangi2Email(templatePath, emailDataObject)

	emailMsg := CreateEmailMessage(sendDueTime, to, from, subject, emailBody, msgType, debugInfo, attachments...)
	emailMsg.Enqueue(emailContext, scheduleSendMailNow, sendGridAdvancedSuppressionManagerGroup)
}

func EnqueueEmail_AndSendNow(
	sendGridAdvancedSuppressionManagerGroup int,
	viewData iEmailView,
	emailContext *EmailContext,
	to []*EmailRecipient,
	from *EmailRecipient,
	subject string,
	msgType IEmailMessageType,
	debugInfo string,
	attachments ...*EmailAttachment) {
	EnqueueEmail(sendGridAdvancedSuppressionManagerGroup, viewData, emailContext, true, time.Now(), to, from, subject, msgType, debugInfo, attachments...)
}
