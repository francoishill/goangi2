package emailUtils

import (
	"crypto/tls"
	"net"
	"net/smtp"
	"time"

	. "github.com/francoishill/goangi2/utils/loggingUtils"
)

const (
	EMAIL_PROVIDER_TYPE__DEFAULT  = 1 << iota
	EMAIL_PROVIDER_TYPE__SENDGRID = 1 << iota
)

var DefaultEmailContext *EmailContext

func CreateEmailContext_DefaultEmailProvider(logger ILogger, authUsername, authPassword, mailHostAndPort string, queueSendingIntervalMinutes int, doNotReplyFrom, adminRecipient, supportRecipient *EmailRecipient) *EmailContext {
	return &EmailContext{
		Logger:                      logger,
		emailProviderType:           EMAIL_PROVIDER_TYPE__DEFAULT,
		doNotReplyFrom:              doNotReplyFrom,
		adminRecipient:              adminRecipient,
		supportRecipient:            supportRecipient,
		authUsername:                authUsername,
		authPassword:                authPassword,
		mailHostAndPort:             mailHostAndPort,
		queueSendingIntervalMinutes: queueSendingIntervalMinutes,
		sendgridApiUser:             "", //Not used for Default email provider
		sendgridApiKey:              "", //Not used for Default email provider
		sendgridTemplateId:          "", //Not used for Default email provider
	}
}

func CreateEmailContext_SendGridProvider(logger ILogger, sendgridApiUser, sendgridApiKey, sendgridTemplateId string, doNotReplyFrom, adminRecipient, supportRecipient *EmailRecipient) *EmailContext {
	return &EmailContext{
		Logger:                      logger,
		emailProviderType:           EMAIL_PROVIDER_TYPE__SENDGRID,
		doNotReplyFrom:              doNotReplyFrom,
		adminRecipient:              adminRecipient,
		supportRecipient:            supportRecipient,
		authUsername:                "", //Not used for SendGrid email provider
		authPassword:                "", //Not used for SendGrid email provider
		mailHostAndPort:             "", //Not used for SendGrid email provider
		queueSendingIntervalMinutes: 0,  //Not used for SendGrid email provider
		sendgridApiUser:             sendgridApiUser,
		sendgridApiKey:              sendgridApiKey,
		sendgridTemplateId:          sendgridTemplateId,
	}
}

type EmailContext struct {
	Logger ILogger

	emailProviderType int

	doNotReplyFrom   *EmailRecipient //A do-not-reply email which is used as the from email of system mail, like Register, Forgot Password, etc
	adminRecipient   *EmailRecipient //Typically the person receiving general Errors and Notices
	supportRecipient *EmailRecipient //Typically the person responsible to respond to payment/transactions issues and failures

	authUsername                string
	authPassword                string
	mailHostAndPort             string
	queueSendingIntervalMinutes int

	sendgridApiUser    string
	sendgridApiKey     string
	sendgridTemplateId string
}

func (this *EmailContext) getSmtpPlainAuth() smtp.Auth {
	host, _, err := net.SplitHostPort(this.mailHostAndPort)
	if err != nil {
		panic(err)
	}
	var auth smtp.Auth
	if len(this.authUsername) > 0 || len(this.authPassword) > 0 {
		auth = smtp.PlainAuth("", this.authUsername, this.authPassword, host)
	} else {
		auth = nil
	}
	return auth
}

func (this *EmailContext) getAuthorizedSmtpClient() *smtp.Client {
	client, err := smtp.Dial(this.mailHostAndPort)
	checkError(err)

	host, _, _ := net.SplitHostPort(this.mailHostAndPort)
	tlsConn := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         host,
	}

	err = client.StartTLS(tlsConn)
	checkError(err)

	auth := this.getSmtpPlainAuth()
	if auth != nil {
		err = client.Auth(auth)
		checkError(err)
	}

	return client
}

func (this *EmailContext) GetDoNotReplyFrom() *EmailRecipient {
	return this.doNotReplyFrom
}

func (this *EmailContext) GetQueueSendingInterval() time.Duration {
	return time.Minute * time.Duration(this.queueSendingIntervalMinutes)
}
