package errors

import (
	"fmt"
	"net/url"
	"strings"

	. "github.com/francoishill/goangi2/utils/debugUtils"
	. "github.com/francoishill/goangi2/utils/loggingUtils"
)

func LogRouterError_And_ExtractUserMessage(loggerToUse ILogger, errorPrefix string, requestUrl *url.URL, remoteAddress string, proxies []string, userAgent string, recoveryObj interface{}) string { //return user message
	logMessage := ""
	userMessage := ""
	mustLogMsg := true
	if expErr, ok := recoveryObj.(*NotLoggedError); ok {
		mustLogMsg = false
		userMessage = expErr.ErrorString
	} else if strMsg, ok := recoveryObj.(string); ok {
		logMessage = fmt.Sprintf(errorPrefix+"'%s'. Request url: '%s'. Remote address: '%s'. Proxies: '%s'. UserAgent: '%s'. Stack trace: '%s'",
			strMsg, requestUrl.String(), remoteAddress, strings.Join(proxies, ","), userAgent, GetFullStackTrace_Pretty())
		userMessage = strMsg
	} else {
		strMsg := fmt.Sprintf("%+v", recoveryObj)
		logMessage = fmt.Sprintf(errorPrefix+"'%s'. Request url: '%s'. Remote address: '%s'. Proxies: '%s'. UserAgent: '%s'. Stack trace: '%s'",
			strMsg, requestUrl.String(), remoteAddress, strings.Join(proxies, ","), userAgent, GetFullStackTrace_Pretty())
		userMessage = strMsg
	}
	if mustLogMsg {
		loggerToUse.Error(logMessage)
	}

	return userMessage
}
