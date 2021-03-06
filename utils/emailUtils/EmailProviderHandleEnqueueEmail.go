package emailUtils

import (
	"github.com/sendgrid/sendgrid-go"
)

func (this *EmailMessage) enqueue_UsingDefault(emailContext *EmailContext, scheduleSendMailNow bool) {
	queuedEmail := CreateQueuedEmail_FromEmailMessage(this)
	queuedEmail.Insert(nil)
	if scheduleSendMailNow {
		ScheduleSendMailNow()
	}
}

func (this *EmailMessage) enqueue_UsingSendGrid(emailContext *EmailContext, scheduleSendMailNow bool, sendGridAdvancedSuppressionManagerGroup int) {
	sendgridClient := sendgrid.NewSendGridClient(emailContext.sendgridApiUser, emailContext.sendgridApiKey)

	message := sendgrid.NewMail()
	for _, to := range this.To {
		if err := message.AddTo(to.EmailAddress); err != nil {
			emailContext.Logger.Error("Unable to AddTo in sendgrid message, error: %s", err.Error())
			continue //Just skip the current user
		}
		message.AddToName(to.FullName)
	}

	if err := message.SetFrom(this.From.EmailAddress); err != nil {
		emailContext.Logger.Error("Unable to SetFrom in sendgrid message, error: %s", err.Error())
	}

	message.SetFromName(this.From.FullName)

	message.SetSubject(this.Subject)
	message.SetHTML(this.Body)

	if emailContext.sendgridTemplateId != "" {
		message.AddFilter("templates", "enabled", "1")
		message.AddFilter("templates", "template_id", emailContext.sendgridTemplateId)
	}

	message.SetSendAt(this.SendDueTime.Unix())

	if sendGridAdvancedSuppressionManagerGroup > 0 {
		message.SetASMGroupID(sendGridAdvancedSuppressionManagerGroup)
	}

	//We handle the 'building' of the sendgrid message synchronously (not on a go routine), but the actual sending we handle on a go subroutine
	go func() {
		defer func() {
			if r := recover(); r != nil {
				emailContext.Logger.Error("UNEXPECTED error occurred in go routine to send SendGrid email, error: %+v", r)
			}
		}()

		maxRetries := 5 //If the max retries is 5, we will actually attempt (6 times) 1x more because we RETRY 5 times after the first fail
		currentRetries := 0
		for currentRetries <= maxRetries {
			currentRetries++

			err := sendgridClient.Send(message)
			if err == nil {
				break
			}
			emailContext.Logger.Error("Unable to send email using sendgrid (will now retry attempt %d), error: %s", currentRetries, err.Error())
		}
	}()
}
