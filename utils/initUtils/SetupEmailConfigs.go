package initUtils

import (
	"fmt"
	. "github.com/francoishill/goangi2/utils/configUtils"
	. "github.com/francoishill/goangi2/utils/emailUtils"
	. "github.com/francoishill/goangi2/utils/loggingUtils"
	"strings"
)

func SetupEmailConfigs(configProvider IConfigContainer, logger ILogger) *EmailContext {
	emailProvider := strings.ToLower(configProvider.DefaultString("email::email_provider", "default"))

	doNotReplyFrom := CreateEmailRecipient(
		configProvider.MustString("email_froms::do_not_reply__full_name"),
		configProvider.MustString("email_froms::do_not_reply__email"),
	)
	adminRecipient := CreateEmailRecipient(
		configProvider.MustString("email_recipients::admin__full_name"),
		configProvider.MustString("email_recipients::admin__email"),
	)
	supportRecipient := CreateEmailRecipient(
		configProvider.MustString("email_recipients::support__full_name"),
		configProvider.MustString("email_recipients::support__email"),
	)

	//TODO: We must probably expand this emailProvider to use the Visitor pattern.
	switch emailProvider {
	case "default":
		authUsername := configProvider.DefaultString("email::auth_username", "")
		authPassword := configProvider.DefaultString("email::auth_password", "")
		mailHostAndPort := configProvider.MustString("email::host_and_port")
		queueSendingIntervalMinutes := configProvider.MustInt("email::queue_sending_interval_minutes")

		DefaultEmailContext = CreateEmailContext_DefaultEmailProvider(logger, authUsername, authPassword, mailHostAndPort, queueSendingIntervalMinutes, doNotReplyFrom, adminRecipient, supportRecipient)
		go StartContinualEmailQueueSender(DefaultEmailContext)
		break
	case "sendgrid":
		sendgridApiUser := configProvider.MustString("email::sendgrid_api_user")
		sendgridApiKey := configProvider.MustString("email::sendgrid_api_key")
		sendgridTemplateId := configProvider.DefaultString("email::sendgrid_template_id", "")
		DefaultEmailContext = CreateEmailContext_SendGridProvider(logger, sendgridApiUser, sendgridApiKey, sendgridTemplateId, doNotReplyFrom, adminRecipient, supportRecipient)
		break
	default:
		panic(fmt.Sprintf("Unsupported Email Provider '%s'", emailProvider))
	}

	logger.Informational("Using Email Provider '%s' for sending emails.", emailProvider)

	return DefaultEmailContext
}
