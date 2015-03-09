package imageUtils

import (
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/nfnt/resize"

	. "github.com/francoishill/goangi2/utils/errorUtils"
)

func localPanicIfHasError(potentialError error) {
	if potentialError != nil {
		panic(potentialError)
	}
}

func localCopyFile(src, dst string, autoCreateParentFolder bool) (int64, error) {
	sf, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer sf.Close()

	if autoCreateParentFolder {
		parentDirPath := path.Dir(strings.Replace(dst, "\\", "/", -1))
		_, err = os.Stat(parentDirPath)
		if err != nil && os.IsNotExist(err) {
			err = os.MkdirAll(parentDirPath, 0600)
			if err != nil {
				return 0, err
			}
		}
	}

	df, err := os.Create(dst)
	if err != nil {
		return 0, err
	}

	defer df.Close()
	return io.Copy(df, sf)
}

func localPathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func ResizeImageFile(sourceFilePath, destinationFilePath string, sizeLimitMegaBytes int64, maxWidth uint, alwaysRemoveSourceFile bool) {
	if alwaysRemoveSourceFile {
		defer os.Remove(sourceFilePath)
	}

	imageFile, err := os.Open(sourceFilePath)
	localPanicIfHasError(err)
	defer imageFile.Close()

	infoOfImageFile, err := imageFile.Stat()
	localPanicIfHasError(err)
	if infoOfImageFile.Size() > sizeLimitMegaBytes*1024*1024 {
		PanicClientError("We are sorry, but this image exceeds the maximum upload size of %d MB", sizeLimitMegaBytes)
	}

	imageObj, err := jpeg.Decode(imageFile)

	failedToReadAsImage := false
	if err != nil && (strings.HasPrefix(err.Error(), jpeg.FormatError("").Error()) || strings.HasPrefix(err.Error(), jpeg.UnsupportedError("").Error())) {
		imageObj, err = png.Decode(imageFile)
		if err != nil && (strings.HasPrefix(err.Error(), png.FormatError("").Error()) || strings.HasPrefix(err.Error(), png.UnsupportedError("").Error())) {
			imageObj, err = gif.Decode(imageFile)
			if err != nil {
				//panic("Image is not a valid JPEG, PNG or GIF")
				failedToReadAsImage = true
			}
		} else {
			localPanicIfHasError(err)
		}
	} else {
		localPanicIfHasError(err)
	}

	if !failedToReadAsImage {
		maxPixelCount := maxWidth * maxWidth

		size := imageObj.Bounds().Size()
		actualPixelCount := uint(size.X * size.Y)
		if actualPixelCount > maxPixelCount {
			resizedImage := resize.Resize(maxWidth, 0, imageObj, resize.Lanczos3)

			folderOfDestFile := filepath.Dir(destinationFilePath)
			dirExists, err := localPathExists(folderOfDestFile)
			localPanicIfHasError(err)
			if !dirExists {
				err = os.MkdirAll(folderOfDestFile, 0600)
				localPanicIfHasError(err)
			}

			outputFile, err := os.Create(destinationFilePath)
			localPanicIfHasError(err)
			defer outputFile.Close()

			err = jpeg.Encode(outputFile, resizedImage, nil)
			localPanicIfHasError(err)
			return
		}
	}

	//No need to resize the image because it does not exceed max pixel count, so just copy it
	_, err = localCopyFile(sourceFilePath, destinationFilePath, true)
	localPanicIfHasError(err)
}
