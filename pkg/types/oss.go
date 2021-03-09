package types

import (
	"fmt"
	"net/url"
	"path"

	oss "github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

var (
	OSSBucketEndpoint string

	ossClient *oss.Client
	ossBucket *oss.Bucket
)

func InitAliyunClients(ctx *cli.Context) error {
	var err error
	var info *oss.BucketInfo
	if ossClient == nil {
		if ossClient, err = oss.New(
			getEndpoint(),
			accessKeyID,
			accessSecret,
			oss.Timeout(120, 1200),
		); err != nil {
			return err
		}
	}

	if ossClient != nil && ossBucket == nil && ossBucketName != "" {
		infoResp, err := ossClient.GetBucketInfo(ossBucketName)
		if err != nil {
			logrus.Warnf("failed to get bucket %s from oss, %v", ossBucketName, err)
			return nil
		}
		info = &infoResp.BucketInfo
		logrus.Infof("Bucket => Region: %s, Name: %s, ACL: %s.", info.Location, info.Name, info.ACL)
		ossBucket, _ = ossClient.Bucket(ossBucketName)

	} else if ossBucketName == "" {
		logrus.Warnf("oss bucket name is not set.")
	}

	if info != nil {
		OSSBucketEndpoint = getBucketEndpoint()
	}

	return nil
}

func getObjectKey(filename string) string {
	return fmt.Sprintf("%s/%s", ossPathPrefix, filename)
}

func getEndpoint() string {
	return fmt.Sprintf("%s://%s%s.%s", ossScheme, ossRegionPrefix, regionID, aliyunAPIDomain)
}

func getBucketEndpoint() string {
	if ossBucketCnameEndpoint != "" {
		return ossBucketCnameEndpoint
	}
	return fmt.Sprintf("%s://%s.%s%s.%s", ossScheme, ossBucketName, ossRegionPrefix, regionID, aliyunAPIDomain)
}

func signOSSFile(filename string) (string, error) {
	key := fmt.Sprintf("%s/%s", ossPathPrefix, filename)
	return ossBucket.SignURL(key, oss.HTTPGet, ossSignURLExpiresSecond)
}

func IsBucketPrivate(acl string) bool {
	return acl != "private"
}

func GetBucketInfo() (*oss.BucketInfo, error) {
	info, err := ossClient.GetBucketInfo(ossBucketName)
	if err != nil {
		return nil, err
	}
	return &info.BucketInfo, nil
}

func GetOSSFileDownloadURL(filename string, isPublic bool) string {
	if !isPublic {
		fileURL, err := signOSSFile(filename)
		if err != nil {
			logrus.Warnf("failed to sign oss object of %s in bucket %s", getObjectKey(filename), ossBucketName)
		} else {
			return fileURL
		}
	}
	filename = url.PathEscape(filename)
	if ossPathPrefix != "" {
		filename = path.Join(ossPathPrefix, filename)
	}
	return fmt.Sprintf("%s/%s", getBucketEndpoint(), filename)
}

func IsObjectExists(filename string) (bool, error) {
	return ossBucket.IsObjectExist(getObjectKey(filename))
}
