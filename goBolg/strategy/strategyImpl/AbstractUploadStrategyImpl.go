package strategyImpl

import (
	"io"
	"log"
	"mime/multipart"
	"path/filepath"

	"goBolg/utils"
)

type AbstractUploadStrategyImpl struct{}

func (a *AbstractUploadStrategyImpl) UploadFile(file multipart.File, path string) (string, error) {
	md5, err := utils.GetMd5(file)
	if err != nil {
		log.Printf("Error calculating MD5: %v", err)
		return "", err
	}

	// Reset the file reader to the beginning
	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		log.Printf("Error resetting file reader: %v", err)
		return "", err
	}

	extName := filepath.Ext(path)
	fileName := md5 + extName

	filePath := filepath.ToSlash(filepath.Join(path, fileName))
	log.Printf("Checking if file exists: %s", filePath)

	if !a.Exists(filePath) {
		log.Printf("Uploading file: %s to path: %s", fileName, path)
		err := a.Upload(path, fileName, file)
		if err != nil {
			log.Printf("Error uploading file: %v", err)
			return "", err
		}
	}

	url := a.GetFileAccessUrl(filePath)
	log.Printf("Getting file access URL for: %s", filePath)
	log.Printf("File access URL: %s", url)
	return url, nil
}

func (a *AbstractUploadStrategyImpl) UploadFileStream(fileName string, inputStream io.Reader, path string) (string, error) {
	err := a.Upload(path, fileName, inputStream)
	if err != nil {
		log.Printf("Error uploading file stream: %v", err)
		return "", err
	}
	filePath := filepath.ToSlash(filepath.Join(path, fileName))
	url := a.GetFileAccessUrl(filePath)
	log.Printf("File stream uploaded successfully: %s", url)
	return url, nil
}

func (a *AbstractUploadStrategyImpl) Exists(filePath string) bool {
	// This method should be implemented by concrete strategies
	log.Printf("Checking if file exists in abstract strategy: %s", filePath)
	return false
}

func (a *AbstractUploadStrategyImpl) Upload(path, fileName string, inputStream io.Reader) error {
	// This method should be implemented by concrete strategies
	fullPath := filepath.ToSlash(filepath.Join(path, fileName))
	log.Printf("Uploading file in abstract strategy: %s", fullPath)
	return nil
}

func (a *AbstractUploadStrategyImpl) GetFileAccessUrl(filePath string) string {
	// This method should be implemented by concrete strategies
	log.Printf("Generating file access URL in abstract strategy: %s", filePath)
	return ""
}
