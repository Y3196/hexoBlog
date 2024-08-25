package strategyImpl

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"goBolg/enums"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type osscurd struct {
	client       *oss.Client
	EndPoint     string
	AccessKeyId  string
	AccessSecret string
	BucketName   string
}

func Aliyun_Oss_Init(EndPoint, AccessKeyId, AccessSecret, BucketName string) (*osscurd, error) {
	var err error
	var oc osscurd
	oc.client, err = oss.New(EndPoint, AccessKeyId, AccessSecret)
	oc.BucketName = BucketName
	oc.AccessKeyId = AccessKeyId
	oc.AccessSecret = AccessSecret
	oc.EndPoint = EndPoint
	if err != nil {
		return nil, err
	}
	return &oc, nil
}

func (o *osscurd) UploadFile(fileHeader *multipart.FileHeader, pathEnum enums.FilePathEnum) (string, error) {
	log.Println("Starting file upload...")
	ts := time.Now().Unix()
	tsStr := strconv.FormatInt(ts, 10)
	before, after, _ := strings.Cut(fileHeader.Filename, ".")
	newFileName := pathEnum.Path + before + "_" + tsStr + "." + after // usage of bucket_file_name

	// 2. bucket
	exist, err := o.client.IsBucketExist(o.BucketName)
	if err != nil {
		log.Printf("Error checking if bucket exists: %v", err)
		return "", err
	}

	if !exist {
		if err := o.client.CreateBucket(o.BucketName); err != nil {
			log.Printf("Error creating bucket: %v", err)
			return "", err
		}
	}

	bucketIns, err := o.client.Bucket(o.BucketName)
	if err != nil {
		log.Printf("Error getting bucket instance: %v", err)
		return "", err
	}

	// 3. upload
	fd, err := fileHeader.Open()
	if err != nil {
		log.Printf("Error opening file header: %v", err)
		return "", err
	}
	fileBytes, err := ioutil.ReadAll(fd)
	if err != nil {
		log.Printf("Error reading file: %v", err)
		return "", err
	}
	if err = bucketIns.PutObject(newFileName, bytes.NewReader(fileBytes)); err != nil {
		log.Printf("Error putting object in bucket: %v", err)
		return "", err
	} else {
		log.Printf("File uploaded successfully with new name: %s", newFileName)
		//返回key
		return newFileName, nil
	}
}

func (o *osscurd) Aliyun_Oss_GetFileURL(fileKey string) (string, error) {
	log.Printf("Generating file URL for key: %s", fileKey)
	// 2. bucket
	exist, err := o.client.IsBucketExist(o.BucketName)
	if err != nil {
		return "", err
	}

	if !exist {
		if err := o.client.CreateBucket(o.BucketName); err != nil {
			return "", err
		}
	}

	bucketIns, err := o.client.Bucket(o.BucketName)
	if err != nil {
		log.Printf("Error getting bucket instance: %v", err)
		return "", err
	}

	// Check if the object exists before trying to get it
	isExist, err := bucketIns.IsObjectExist(fileKey)
	if err != nil {
		log.Printf("Error checking if object exists: %v", err)
		return "", err
	}
	if !isExist {
		log.Printf("The specified key does not exist: %s", fileKey)
		return "", errors.New("The specified key does not exist.")
	}

	//获取文件url路径
	signedURL, err := bucketIns.SignURL(fileKey, oss.HTTPGet, 60)
	if err != nil {
		log.Printf("Error signing URL: %v", err)
		return "", err
	}
	// 解析URL
	parsedURL, err := url.Parse(signedURL)
	if err != nil {
		log.Printf("Error parsing signed URL: %v", err)
		return "", err
	}

	// 清除查询参数
	parsedURL.RawQuery = ""

	// 重建无查询参数的URL
	cleanedURL := parsedURL.String()
	log.Printf("Generated file URL: %s", cleanedURL)
	return fmt.Sprintf("%s", cleanedURL), nil
}
