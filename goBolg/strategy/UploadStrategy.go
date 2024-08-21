package strategy

import (
	"goBolg/enums"
	"mime/multipart"
)

type UploadStrategy interface {
	UploadFile(fileHeader *multipart.FileHeader, pathEnum enums.FilePathEnum) (string, error)
	Aliyun_Oss_GetFileURL(fileKey string) (string, error)
}
