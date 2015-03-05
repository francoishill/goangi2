package httpUtils

import (
	"io"
	"io/ioutil"
	"mime/multipart"
	"os"
	"path/filepath"
)

func AcceptUploadedFileAndSaveToFolder(filePrefix, folderToSaveTo string, file multipart.File) (string, error) {
	var err error
	err = os.MkdirAll(folderToSaveTo, 0600)
	if err != nil {
		return "", err
	}

	t, _ := ioutil.TempFile(folderToSaveTo, filePrefix)
	defer t.Close()

	_, err = io.Copy(t, file)
	if err != nil {
		return "", err
	}

	fullFilePath, err := filepath.Abs(t.Name())
	if err != nil {
		return "", err
	}

	return fullFilePath, nil
}
