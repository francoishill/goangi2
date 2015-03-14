package initUtils

import (
	. "github.com/francoishill/goangi2/utils/configUtils"
	. "github.com/francoishill/goangi2/utils/emailUtils"
	. "github.com/francoishill/goangi2/utils/loggingUtils"
)

func SetupEmailConfigs(configProvider IConfigContainer, logger ILogger) *EmailContext {
	authUsername := configProvider.DefaultString("email::auth_username", "")
	authPassword := configProvider.DefaultString("email::auth_password", "")
	mailHostAndPort := configProvider.MustString("email::host_and_port")
	queueSendingIntervalMinutes := configProvider.MustInt("email::queue_sending_interval_minutes")
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

	DefaultEmailContext = CreateEmailContext(logger, authUsername, authPassword, mailHostAndPort, queueSendingIntervalMinutes, doNotReplyFrom, adminRecipient, supportRecipient)
	go StartContinualEmailQueueSender(DefaultEmailContext)

	return DefaultEmailContext
}
