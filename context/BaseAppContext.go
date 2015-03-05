package context

import (
	"github.com/francoishill/goangi2/utils/httpUtils"
	"github.com/francoishill/goangi2/utils/imageUtils"
	. "github.com/francoishill/goangi2/utils/loggingUtils"
	"io/ioutil"
	"mime/multipart"
	"os"
	"path/filepath"
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

func (this *BaseAppContext) checkError(err error) {
	if err != nil {
		panic(err)
	}
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

func (this *BaseAppContext) getTempImageFileFullPath(fileNameOnly string) string {
	return filepath.Join(this.UploadDirectory, fileNameOnly)
}

func (this *BaseAppContext) ReadTempImageFileBytes(fileNameOnly string) []byte {
	fullTempFilePath := this.getTempImageFileFullPath(fileNameOnly)
	fileBytes, err := ioutil.ReadFile(fullTempFilePath)
	this.checkError(err)
	return fileBytes
}

func (this *BaseAppContext) UploadAndResizeImageToTempUploadDir(file multipart.File, originalFilenamePrefix, resizedFilenamePrefix string, maxImageWidth uint) string {
	originalTempFile, err := httpUtils.AcceptUploadedFileAndSaveToFolder(originalFilenamePrefix, this.UploadDirectory, file)
	this.checkError(err)

	resizedTempFilePathObj, err := ioutil.TempFile(filepath.Dir(originalTempFile), resizedFilenamePrefix)
	this.checkError(err)
	defer resizedTempFilePathObj.Close()

	resizedTempFilePath, err := filepath.Abs(resizedTempFilePathObj.Name())
	this.checkError(err)

	imageUtils.ResizeImageFile(originalTempFile, resizedTempFilePath, maxImageWidth)
	defer os.Remove(originalTempFile)

	return resizedTempFilePath
}

func (this *BaseAppContext) UploadResizedProfilePic(file multipart.File) string {
	origImageFilenamePrefix := "temp-profilepic-origsize-"
	resizedImageFilenamePrefix := "temp-profilepic-resized-"
	return this.UploadAndResizeImageToTempUploadDir(file, origImageFilenamePrefix, resizedImageFilenamePrefix, this.MaxProfilePicWidth)
}
