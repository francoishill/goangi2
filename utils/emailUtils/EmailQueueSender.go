package emailUtils

import (
	"fmt"
	"strings"
	"time"

	. "github.com/francoishill/goangi2/utils/debugUtils"
)

var sendMailNowTriggered bool = true

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

func onEmailQueueSendError(emailContext *EmailContext) {
	if r := recover(); r != nil {
		emailContext.Logger.Error("Error in trying to execute email queue sender, error: %+v", r)
	}
}

func sendEmail(emailContext *EmailContext, queuedEmail *QueuedEmail) {
	defer func() {
		if r := recover(); r != nil {
			originalError := fmt.Errorf("%+v. Stack trace:\n%s", r, GetFullStackTrace_Pretty())

			func() {
				defer func() {
					if r2 := recover(); r2 != nil {
						emailContext.Logger.Error("Unable to save queued email error: %+v", r2)
					}
				}()
				queuedEmail.Update_SetSendError(nil, originalError.Error())
			}()

			expandedError := fmt.Errorf("{Additional send email debug info: %s} %s", queuedEmail.DebugInfo, originalError.Error())
			emailContext.Logger.Error(expandedError.Error())
			panic(expandedError)
		}
	}()

	client := emailContext.getAuthorizedSmtpClient()

	err := client.Mail(queuedEmail.FromEmail)
	checkError(err)

	toEmails := queuedEmail.GetSplittedRecipientEmailAddresses()
	for _, toEmail := range toEmails {
		err = client.Rcpt(toEmail)
		checkError(err)
	}

	w, err := client.Data()
	checkError(err)

	_, err = w.Write([]byte(queuedEmail.Content))
	checkError(err)

	err = w.Close()
	checkError(err)

	err = client.Quit()
	checkError(err)

	queuedEmail.Update_MarkAsSuccessfullySent(nil)
}

func sendAllMailsInQueue(emailContext *EmailContext) {
	defer onEmailQueueSendError(emailContext)

	emailContext.Logger.Debug("Starting to send all mails")

	caughtErrors := []string{}

	unsentQueuedEmails := (&QueuedEmail{}).List_Unsent(nil, nil)
	for _, queuedEmail := range unsentQueuedEmails {
		if queuedEmail.SendDueTime.After(time.Now()) {
			continue
		}
		func() {
			defer func() {
				if r := recover(); r != nil {
					caughtErrors = append(caughtErrors, fmt.Sprintf("%+v", r))
				}
			}()
			sendEmail(emailContext, queuedEmail)
		}()
	}

	if len(caughtErrors) > 0 {
		panic(fmt.Errorf("Unable to send %d emails. Error list: %s", len(caughtErrors), strings.Join(caughtErrors, "|")))
	}
}

func ScheduleSendMailNow() {
	sendMailNowTriggered = true
}

func StartContinualEmailQueueSender(emailContext *EmailContext) {
	for {
		sendAllMailsInQueue(emailContext)
		expectedWaitDueTime := time.Now().Add(emailContext.GetQueueSendingInterval())
		for {
			time.Sleep(1 * time.Second)

			if sendMailNowTriggered {
				sendMailNowTriggered = false
				break
			}

			if !expectedWaitDueTime.After(time.Now()) {
				break
			}
		}
	}
}
