package context

import (
	. "github.com/francoishill/goangi2/utils/loggingUtils"
	"strings"
)

var DefaultBaseAppContext *BaseAppContext

type BaseAppContext struct {
	Logger               ILogger
	baseAppUrl_WithSlash string
	baseAppUrl_NoSlash   string
	MaxProfilePicWidth   uint
	UploadDirectory      string
	ProfilePicsDirectory string
}

func CreateBaseAppContext(logger ILogger, baseAppUrl string, maxProfilePicWidth uint, uploadDir, profilePicsDir string) *BaseAppContext {
	baseAppUrlNoSlash := strings.TrimRight(baseAppUrl, "/")

	return &BaseAppContext{
		Logger:               logger,
		baseAppUrl_WithSlash: baseAppUrlNoSlash + "/",
		baseAppUrl_NoSlash:   baseAppUrlNoSlash,
		MaxProfilePicWidth:   maxProfilePicWidth,
		UploadDirectory:      uploadDir,
		ProfilePicsDirectory: profilePicsDir,
	}
}

func (this *BaseAppContext) GenerateAppRelativeUrl(partAfterBaseUrl string) string {
	if partAfterBaseUrl == "" {
		return this.baseAppUrl_NoSlash
	}

	if partAfterBaseUrl[0] == '/' {
		return this.baseAppUrl_NoSlash + partAfterBaseUrl
	} else {
		return this.baseAppUrl_WithSlash + partAfterBaseUrl
	}
}
