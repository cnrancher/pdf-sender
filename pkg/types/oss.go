package types

import (
	"fmt"
	"net/url"

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
			Aliyun.AccessKey,
			Aliyun.AccessSecret,
			oss.Timeout(120, 1200),
		); err != nil {
			return err
		}
	}

	if ossClient != nil && ossBucket == nil && Aliyun.OSSBucket != "" {
		infoResp, err := ossClient.GetBucketInfo(Aliyun.OSSBucket)
		if err != nil {
			logrus.Warnf("failed to get bucket %s from oss, %v", Aliyun.OSSBucket, err)
			return nil
		}
		info = &infoResp.BucketInfo
		logrus.Infof("Bucket => Region: %s, Name: %s, ACL: %s.", info.Location, info.Name, info.ACL)
		ossBucket, _ = ossClient.Bucket(Aliyun.OSSBucket)

	} else if Aliyun.OSSBucket == "" {
		logrus.Warnf("oss bucket name is not set.")
	}

	if info != nil {
		OSSBucketEndpoint = getBucketEndpoint()
	}

	return nil
}

func getObjectKey(prefix, filename string) string {
	key := filename
	if prefix == "" {
		prefix = Aliyun.OSSPathPrefix
	}
	if prefix != "" {
		key = prefix + "/" + key
	}
	return key
}

func getEndpoint() string {
	return fmt.Sprintf("%s://%s%s.%s", Aliyun.OSSScheme, ossRegionPrefix, Aliyun.Region, aliyunAPIDomain)
}

func getBucketEndpoint() string {
	if Aliyun.OSSCnameEndpoint != "" {
		return Aliyun.OSSCnameEndpoint
	}
	return fmt.Sprintf("%s://%s.%s%s.%s", Aliyun.OSSScheme, Aliyun.OSSBucket, ossRegionPrefix, Aliyun.Region, aliyunAPIDomain)
}

func signOSSFile(prefix, filename string) (string, error) {
	return ossBucket.SignURL(getObjectKey(prefix, filename), oss.HTTPGet, Aliyun.OSSSignURLExpiresSecond)
}

func IsBucketPrivate(acl string) bool {
	return acl != "private"
}

func GetBucketInfo() (*oss.BucketInfo, error) {
	info, err := ossClient.GetBucketInfo(Aliyun.OSSBucket)
	if err != nil {
		return nil, err
	}
	return &info.BucketInfo, nil
}

func GetOSSFileDownloadURL(filename string, prefix string, isPublic bool) string {
	if !isPublic {
		fileURL, err := signOSSFile(prefix, filename)
		if err != nil {
			logrus.Warnf("failed to sign oss object of %s in bucket %s", getObjectKey(prefix, filename), Aliyun.OSSBucket)
		} else {
			return fileURL
		}
	}
	return fmt.Sprintf("%s/%s", getBucketEndpoint(), getObjectKey(prefix, url.PathEscape(filename)))
}

func IsObjectExists(prefixOverride, filename string) (bool, error) {
	return ossBucket.IsObjectExist(getObjectKey(prefixOverride, filename))
}
