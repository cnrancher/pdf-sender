package types

import (
	"github.com/urfave/cli"
)

const (
	ossRegionPrefix = "oss-"
	aliyunAPIDomain = "aliyuncs.com"
)

var (
	ossSignURLExpiresSecond int64

	ossScheme,
	ossBucketCnameEndpoint,
	ossBucketName,
	ossPathPrefix,
	regionID,
	accessKeyID,
	accessSecret,
	signName,
	templateCode string
)

func GetFlags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:        "ali-oss-endpoint-scheme",
			EnvVar:      "ALI_OSS_ENDPOINT_SCHEME",
			Value:       "https",
			Destination: &ossScheme,
		},
		cli.StringFlag{
			Name:        "ali-oss-path-prefix",
			EnvVar:      "ALI_OSS_PATH_PREFIX",
			Usage:       "Path prefix for documents in oss bucket.",
			Destination: &ossPathPrefix,
		},
		cli.StringFlag{
			Name:        "ali-oss-bucket-name",
			EnvVar:      "ALI_OSS_BUCKET_NAME",
			Usage:       "The private oss bucket name.",
			Destination: &ossBucketName,
		},
		cli.Int64Flag{
			Name:        "ali-oss-sign-expire-second",
			EnvVar:      "ALI_OSS_SIGN_EXPIRE_SECOND",
			Usage:       "The expired second of signing the private oss object.",
			Destination: &ossSignURLExpiresSecond,
			Value:       86400,
		},
		cli.StringFlag{
			Name:        "ali-oss-cname-endpoint",
			EnvVar:      "ALI_OSS_CNAME_ENDPOINT",
			Usage:       "The CNAME endpoint of the bucket when generating file download URL.",
			Destination: &ossBucketCnameEndpoint,
		},
		cli.StringFlag{
			Name:        "ali-region-id",
			EnvVar:      "ALI_REGION_ID",
			Usage:       "The region id of accessing aliyun API, will also use oss bucket region.",
			Destination: &regionID,
			Required:    true,
		},
		cli.StringFlag{
			Name:        "ali-access-key-id",
			EnvVar:      "ALI_ACCESS_KEYID",
			Usage:       "The access key of accessing aliyun API.",
			Destination: &accessKeyID,
			Required:    true,
		},
		cli.StringFlag{
			Name:        "ali-access-secret",
			EnvVar:      "ALI_ACCESS_SECRET",
			Usage:       "The access secret of accessing aliyun API",
			Destination: &accessSecret,
			Required:    true,
		},
		cli.StringFlag{
			Name:        "ali-sign-name",
			EnvVar:      "ALI_SIGN_NAME",
			Usage:       "The signature of sms",
			Destination: &signName,
		},
		cli.StringFlag{
			Name:        "ali-template-code",
			EnvVar:      "ALI_TEMPLATE_CODE",
			Usage:       "The template code of aliyun sms",
			Destination: &templateCode,
		},
	}
}
