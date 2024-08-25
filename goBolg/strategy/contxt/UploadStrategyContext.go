package contxt

import (
	"fmt"
	"goBolg/config"
	"goBolg/enums" // Ensure enums package is imported
	"goBolg/strategy"
	"goBolg/strategy/strategyImpl"
	"log"
	"mime/multipart"
)

type UploadStrategyContext struct {
	StrategyMap map[string]strategy.UploadStrategy
	Config      *config.AppConfig
}

func NewUploadStrategyContext(config *config.AppConfig) (*UploadStrategyContext, error) {
	log.Printf("Loaded Upload Config: %+v", config.Upload)
	if config == nil {
		log.Fatalf("AppConfig is nil")
	}
	if config.Upload == nil {
		log.Fatalf("UploadConfig is nil")
	}
	if config.Upload.Oss.Endpoint == "" || config.Upload.Oss.AccessKeyId == "" || config.Upload.Oss.AccessKeySecret == "" || config.Upload.Oss.BucketName == "" {
		log.Fatalf("OSS configuration is incomplete")
	}
	ossStrategy, err := strategyImpl.Aliyun_Oss_Init(config.Upload.Oss.Endpoint, config.Upload.Oss.AccessKeyId, config.Upload.Oss.AccessKeySecret, config.Upload.Oss.BucketName)
	if err != nil {
		log.Printf("Error initializing OSS strategy: %v", err)
		return nil, err
	}
	log.Printf("OSS strategy initialized successfully")

	return &UploadStrategyContext{
		StrategyMap: map[string]strategy.UploadStrategy{
			"oss": ossStrategy,
		},
		Config: config,
	}, nil
}

func (c *UploadStrategyContext) ExecuteUploadStrategy(fileHeader *multipart.FileHeader, pathEnum enums.FilePathEnum) (string, error) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Recovered from panic in ExecuteUploadStrategy: %v", r)
		}
	}()

	log.Printf("ExecuteUploadStrategy: Starting upload...")

	if c == nil {
		log.Fatalf("UploadStrategyContext is nil!")
		return "", fmt.Errorf("upload strategy context is not initialized")
	}

	if c.Config == nil {
		log.Fatalf("UploadStrategyContext: c.Config is nil!")
		return "", fmt.Errorf("upload strategy configuration is not initialized")
	}

	if c.Config.Upload == nil {
		log.Fatalf("UploadStrategyContext: c.Config.Upload is nil!")
		return "", fmt.Errorf("upload configuration is not initialized")
	}

	log.Printf("Executing upload strategy for mode: %s", c.Config.Upload.Mode)
	strategy, ok := c.StrategyMap[c.Config.Upload.Mode]
	if !ok {
		err := fmt.Errorf("unsupported upload mode: %s", c.Config.Upload.Mode)
		log.Printf("Error: %v", err)
		return "", err
	}

	fileKey, err := strategy.UploadFile(fileHeader, pathEnum)
	if err != nil {
		log.Printf("Error uploading file: %v", err)
		return "", err
	}

	fileUrl, err := strategy.Aliyun_Oss_GetFileURL(fileKey)
	if err != nil {
		log.Printf("Error getting file URL: %v", err)
		return "", err
	}

	log.Printf("File uploaded successfully: %s", fileUrl)
	return fileUrl, nil
}
