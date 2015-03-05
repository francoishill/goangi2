package context

import (
	"fmt"
	"github.com/francoishill/goangi2/utils/httpUtils"
	"github.com/francoishill/goangi2/utils/imageUtils"
	. "github.com/francoishill/goangi2/utils/loggingUtils"
	. "github.com/francoishill/goangi2/utils/osUtils"
	"io/ioutil"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
)

var DefaultBaseAppContext *BaseAppContext

type BaseAppContext struct {
	Logger                  ILogger
	baseAppUrl_WithSlash    string
	baseAppUrl_NoSlash      string
	MaxProfilePicWidth      uint
	UploadDirectory         string
	ProfilePicsDirectory    string
	UploadedImagesDirectory string
}

func (this *BaseAppContext) checkError(err error) {
	if err != nil {
		panic(err)
	}
}

func CreateBaseAppContext(logger ILogger, baseAppUrl string, maxProfilePicWidth uint, uploadDir, profilePicsDir, uploadedImagesDir string) *BaseAppContext {
	baseAppUrlNoSlash := strings.TrimRight(baseAppUrl, "/")

	if !DirectoryExists(uploadDir) {
		panic("Uploads directory does not exist: " + uploadDir)
	}

	if !DirectoryExists(profilePicsDir) {
		panic("Profile pics directory does not exist: " + profilePicsDir)
	}

	if !DirectoryExists(uploadedImagesDir) {
		panic("Upload images directory does not exist: " + uploadedImagesDir)
	}

	return &BaseAppContext{
		Logger:                  logger,
		baseAppUrl_WithSlash:    baseAppUrlNoSlash + "/",
		baseAppUrl_NoSlash:      baseAppUrlNoSlash,
		MaxProfilePicWidth:      maxProfilePicWidth,
		UploadDirectory:         uploadDir,
		ProfilePicsDirectory:    profilePicsDir,
		UploadedImagesDirectory: uploadedImagesDir,
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

func (this *BaseAppContext) getProfilePicFileFullPath(userId int64) string {
	return filepath.Join(this.ProfilePicsDirectory, fmt.Sprintf("%d", userId))
}

func (this *BaseAppContext) getUploadedImagePermanentFullPath(imageFileName string) string {
	return filepath.Join(this.UploadedImagesDirectory, fmt.Sprintf("%s", imageFileName))
}

func (this *BaseAppContext) ReadTempImageFileBytes(fileNameOnly string) []byte {
	fullTempFilePath := this.getTempImageFileFullPath(fileNameOnly)
	fileBytes, err := ioutil.ReadFile(fullTempFilePath)
	this.checkError(err)
	return fileBytes
}

func (this *BaseAppContext) ReadPermanentImageFileBytes(fileNameOnly string) []byte {
	fullTempFilePath := this.getUploadedImagePermanentFullPath(fileNameOnly)
	fileBytes, err := ioutil.ReadFile(fullTempFilePath)
	this.checkError(err)
	return fileBytes
}

func (this *BaseAppContext) DeleteTempImageFile(fileNameOnly string) {
	fullTempFilePath := this.getTempImageFileFullPath(fileNameOnly)
	err := os.Remove(fullTempFilePath)
	this.checkError(err)
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

func (this *BaseAppContext) MoveTempProfilePicToPermanentFolder(profilePicFileNameOnly string, userId int64) {
	origTempFullFilePath := this.getTempImageFileFullPath(profilePicFileNameOnly)
	newPermanentFullFilePath := this.getProfilePicFileFullPath(userId)

	err := os.Rename(origTempFullFilePath, newPermanentFullFilePath)
	this.checkError(err)
}

func (this *BaseAppContext) MoveTempImageFileToPermanentFolder(tempFileNameOnly, finalImageName string) {
	origTempFullFilePath := this.getTempImageFileFullPath(tempFileNameOnly)
	newPermanentFullFilePath := this.getUploadedImagePermanentFullPath(finalImageName)

	err := os.Rename(origTempFullFilePath, newPermanentFullFilePath)
	this.checkError(err)
}
