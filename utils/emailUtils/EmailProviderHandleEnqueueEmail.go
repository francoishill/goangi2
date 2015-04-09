package emailUtils

import (
	"github.com/francoishill/sendgrid-go"
)

func (this *EmailMessage) enqueue_UsingDefault(emailContext *EmailContext, scheduleSendMailNow bool) {
	queuedEmail := CreateQueuedEmail_FromEmailMessage(this)
	queuedEmail.Insert(nil)
	if scheduleSendMailNow {
		ScheduleSendMailNow()
	}
}

func (this *EmailMessage) enqueue_UsingSendGrid(emailContext *EmailContext, scheduleSendMailNow bool) {
	sendgridClient := sendgrid.NewSendGridClient(emailContext.sendgridApiUser, emailContext.sendgridApiKey)

	message := sendgrid.NewMail()
	for _, to := range this.To {
		if err := message.AddTo(to.EmailAddress); err != nil {
			panic(err)
		}
		message.AddToName(to.FullName)
	}

	if err := message.SetFrom(this.From.EmailAddress); err != nil {
		panic(err)
	}

	message.SetFromName(this.From.FullName)

	message.SetSubject(this.Subject)
	message.SetHTML(this.Body)

	if emailContext.sendgridTemplateId != "" {
		message.AddFilter("templates", "enabled", "1")
		message.AddFilter("templates", "template_id", emailContext.sendgridTemplateId)
	}

	message.SetSendAtTime(this.SendDueTime)

	//We handle the 'building' of the sendgrid message synchronously (not on a go routine), but the actual sending we handle on a go subroutine
	go func() {
		if err := sendgridClient.Send(message); err != nil {
			panic(err)
		}
	}()
}
