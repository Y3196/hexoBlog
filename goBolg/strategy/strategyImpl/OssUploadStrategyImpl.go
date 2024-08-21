package strategyImpl

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"goBolg/enums"
	"io/ioutil"
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
	ts := time.Now().Unix()
	tsStr := strconv.FormatInt(ts, 10)
	before, after, _ := strings.Cut(fileHeader.Filename, ".")
	newFileName := pathEnum.Path + before + "_" + tsStr + "." + after // usage of bucket_file_name

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
		return "", err
	}

	// 3. upload
	fd, err := fileHeader.Open()
	if err != nil {
		return "", err
	}
	fileBytes, err := ioutil.ReadAll(fd)
	if err != nil {
		return "", err
	}
	if err = bucketIns.PutObject(newFileName, bytes.NewReader(fileBytes)); err != nil {
		return "", err
	} else {
		//返回key
		return newFileName, nil
	}
}

func (o *osscurd) Aliyun_Oss_GetFileURL(fileKey string) (string, error) {
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
		return "", err
	}

	// Check if the object exists before trying to get it
	isExist, err := bucketIns.IsObjectExist(fileKey)
	if err != nil {
		return "", err
	}
	if !isExist {
		return "", errors.New("The specified key does not exist.")
	}

	//获取文件url路径
	signedURL, err := bucketIns.SignURL(fileKey, oss.HTTPGet, 60)
	if err != nil {
		return "", err
	}
	// 解析URL
	parsedURL, err := url.Parse(signedURL)
	if err != nil {
		return "", err
	}

	// 清除查询参数
	parsedURL.RawQuery = ""

	// 重建无查询参数的URL
	cleanedURL := parsedURL.String()
	return fmt.Sprintf("%s", cleanedURL), nil
}
