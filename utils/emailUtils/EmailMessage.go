package emailUtils

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	. "github.com/francoishill/goangi2/utils/bufferUtils"
)

func CreateEmailMessage(sendDueTime time.Time, to []*EmailRecipient, from *EmailRecipient, subject, body string, msgType IEmailMessageType, debugInfo string, attachments ...*EmailAttachment) *EmailMessage {
	if len(to) == 0 {
		panic(fmt.Errorf("Please specify at least ONE To address"))
	}

	return &EmailMessage{
		To:          to,
		From:        from,
		Subject:     subject,
		Body:        body,
		Type:        msgType,
		Attachments: attachments,
		DebugInfo:   debugInfo,
		SendDueTime: sendDueTime,
	}
}

type EmailMessage struct {
	To          []*EmailRecipient
	From        *EmailRecipient
	Subject     string
	Body        string
	Type        IEmailMessageType
	Attachments []*EmailAttachment
	SendDueTime time.Time
	DebugInfo   string
}

func (this *EmailMessage) getContentTypeString() string {
	visitor := new(getContentTypeForEmailMessageTypeVisitor)
	this.Type.Accept(visitor)
	return visitor.ContentTypeString
}

func (this *EmailMessage) getFormattedToAddressList(onlyEmailAddress bool) []string {
	formattedToAddresses := []string{}
	for _, to := range this.To {
		if onlyEmailAddress {
			formattedToAddresses = append(formattedToAddresses, to.EmailAddress)
		} else {
			formattedToAddresses = append(formattedToAddresses, to.formatFullnameAndAddress())
		}
	}
	return formattedToAddresses
}

func (this *EmailMessage) getCommaSeparatedFormattedToAddresses(onlyEmailAddress bool) string {
	formattedToAddresses := this.getFormattedToAddressList(onlyEmailAddress)
	return strings.Join(formattedToAddresses, ",")
}

func (this *EmailMessage) GetContent() string {
	buf := CreateByteBufferWithPanic(bytes.NewBuffer(nil))

	onlyEmailAddress := false
	formattedToAddresses := this.getFormattedToAddressList(onlyEmailAddress)
	toAddressesFormattedAsContent := strings.Join(formattedToAddresses, ",")

	// create mail content
	buf.WriteString_PanicOnError("From: " + this.From.formatFullnameAndAddress() + "\r\n")
	buf.WriteString_PanicOnError("To: " + toAddressesFormattedAsContent + "\r\n")
	buf.WriteString_PanicOnError("Subject: " + this.Subject + "\r\n")
	buf.WriteString_PanicOnError("MIME-Version: 1.0\r\n")

	boundary := "f46d043c813270fc6b04c2d223da"

	if len(this.Attachments) > 0 {
		buf.WriteString_PanicOnError("Content-Type: multipart/mixed; boundary=" + boundary + "\r\n")
		buf.WriteString_PanicOnError("--" + boundary + "\r\n")
	}

	buf.WriteString_PanicOnError("Content-Type: " + this.getContentTypeString() + "\r\n\r\n")
	buf.WriteString_PanicOnError(this.Body)

	//Attachment code example from: https://github.com/scorredoira/email/blob/master/email.go
	if len(this.Attachments) > 0 {
		for _, attachment := range this.Attachments {
			buf.WriteString_PanicOnError("\r\n\r\n--" + boundary + "\r\n")

			if attachment.Inline {
				buf.WriteString_PanicOnError("Content-Type: message/rfc822\r\n")
				buf.WriteString_PanicOnError("Content-Disposition: inline; filename=\"" + attachment.Filename + "\"\r\n\r\n")

				buf.Write_PanicOnError(attachment.Data)
			} else {
				buf.WriteString_PanicOnError("Content-Type: application/octet-stream\r\n")
				buf.WriteString_PanicOnError("Content-Transfer-Encoding: base64\r\n")
				buf.WriteString_PanicOnError("Content-Disposition: attachment; filename=\"" + attachment.Filename + "\"\r\n\r\n")

				b := make([]byte, base64.StdEncoding.EncodedLen(len(attachment.Data)))
				//TODO: This attachment.Data can become quite large, especially if we attach large files, this will be stored in the Email Queue DB table. Is this the correct + best way?
				base64.StdEncoding.Encode(b, attachment.Data)
				buf.Write_PanicOnError(b)
			}

			buf.WriteString_PanicOnError("\r\n--" + boundary)
		}

		buf.WriteString_PanicOnError("--")
	}

	return buf.GetString()
}

func (this *EmailMessage) Enqueue(emailContext *EmailContext, scheduleSendMailNow bool) {
	switch emailContext.emailProviderType {
	case EMAIL_PROVIDER_TYPE__DEFAULT:
		this.enqueue_UsingDefault(emailContext, scheduleSendMailNow)
		break
	case EMAIL_PROVIDER_TYPE__SENDGRID:
		this.enqueue_UsingSendGrid(emailContext, scheduleSendMailNow)
		break
	default:
		panic(fmt.Sprintf("Cannot enqueue email, unsupported email provider type %d", emailContext.emailProviderType))
		break
	}
}
