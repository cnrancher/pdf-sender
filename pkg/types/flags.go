package types

import (
	"github.com/urfave/cli"
)

const (
	ossRegionPrefix = "oss-"
	aliyunAPIDomain = "aliyuncs.com"
)

func GetFlags() []cli.Flag {
	return []cli.Flag{
		cli.IntFlag{
			Name:        "port,p",
			EnvVar:      "HTTP_PORT",
			Value:       8080,
			Destination: &Config.Port,
		},
		cli.StringFlag{
			Name:        "smtp-user",
			EnvVar:      "SMTP_USER",
			Usage:       "The smtp server user for sending emails.",
			Destination: &Email.User,
		},
		cli.StringFlag{
			Name:        "smtp-pwd",
			EnvVar:      "SMTP_PWD",
			Usage:       "the smtp user password for sending emails.",
			Destination: &Email.Password,
		},
		cli.StringFlag{
			Name:        "smtp-endpoint",
			EnvVar:      "SMTP_ENDPOINT",
			Usage:       "The smtp server address or domain name.",
			Destination: &Email.Endpoint,
		},
		cli.IntFlag{
			Name:        "smtp-port",
			EnvVar:      "SMTP_PORT",
			Usage:       "The smtp server port",
			Destination: &Email.Port,
		},
		cli.StringSliceFlag{
			Name:   "smtp-rancher-to-day",
			EnvVar: "SMTP_RANCHER_TO_DAY",
			Usage:  "The email receivers of the pdf-sender daily cron information collector.",
		},
		cli.StringSliceFlag{
			Name:   "smtp-rancher-to-mon",
			EnvVar: "SMTP_RANCHER_TO_MON",
			Usage:  "The email receivers of the pdf-sender monthly cron information collector.",
		},
		cli.StringFlag{
			Name:        "smtp-daily-cron,d",
			EnvVar:      "DAY_CRON",
			Usage:       "The daily cron string for sending statistics.",
			Destination: &Email.CRONDaily,
			Value:       "30 1 * * ?",
		},
		cli.StringFlag{
			Name:        "stmp-monthly-cron,m",
			EnvVar:      "MON_CRON",
			Usage:       "The monthly cron string for sending statistics.",
			Destination: &Email.CRONMonthly,
			Value:       "0 2 1 * ?",
		},
		cli.StringFlag{
			Name:        "smtp-sender-email,s",
			Usage:       "The sender email address.",
			EnvVar:      "SENDER_EMAIL",
			Destination: &Email.Sender,
		},

		cli.StringFlag{
			Name:        "db-host-ip",
			EnvVar:      "DB_HOST_IP",
			Usage:       "The backend mysql hostname or IP.",
			Destination: &DB.HostIP,
		},
		cli.IntFlag{
			Name:        "db-port",
			EnvVar:      "DB_PORT",
			Usage:       "The backend mysql port.",
			Destination: &DB.Port,
			Value:       3306,
		},
		cli.StringFlag{
			Name:        "db-name",
			EnvVar:      "DB_NAME",
			Usage:       "The db name of mysql backend.",
			Destination: &DB.Name,
			Value:       "pdf",
		},
		cli.StringFlag{
			Name:        "db-username",
			EnvVar:      "DB_USERNAME",
			Usage:       "The username of mysql backend.",
			Destination: &DB.Username,
		},
		cli.StringFlag{
			Name:        "db-password",
			EnvVar:      "DB_PASSWORD",
			Usage:       "The password of mysql user.",
			Destination: &DB.Password,
		},
		cli.StringFlag{
			Name:        "ali-oss-endpoint-scheme",
			EnvVar:      "ALI_OSS_ENDPOINT_SCHEME",
			Value:       "https",
			Destination: &Aliyun.OSSScheme,
		},
		cli.StringFlag{
			Name:        "ali-oss-path-prefix",
			EnvVar:      "ALI_OSS_PATH_PREFIX",
			Usage:       "Path prefix for documents in oss bucket.",
			Destination: &Aliyun.OSSPathPrefix,
		},
		cli.StringFlag{
			Name:        "ali-oss-bucket-name",
			EnvVar:      "ALI_OSS_BUCKET_NAME",
			Usage:       "The private oss bucket name.",
			Destination: &Aliyun.OSSBucket,
		},
		cli.Int64Flag{
			Name:        "ali-oss-sign-expire-second",
			EnvVar:      "ALI_OSS_SIGN_EXPIRE_SECOND",
			Usage:       "The expired second of signing the private oss object.",
			Destination: &Aliyun.OSSSignURLExpiresSecond,
			Value:       86400,
		},
		cli.StringFlag{
			Name:        "ali-oss-cname-endpoint",
			EnvVar:      "ALI_OSS_CNAME_ENDPOINT",
			Usage:       "The CNAME endpoint of the bucket when generating file download URL.",
			Destination: &Aliyun.OSSCnameEndpoint,
		},
		cli.StringFlag{
			Name:        "ali-region-id",
			EnvVar:      "ALI_REGION_ID",
			Usage:       "The region id of accessing aliyun API, will also use oss bucket region.",
			Destination: &Aliyun.Region,
		},
		cli.StringFlag{
			Name:        "ali-access-key-id",
			EnvVar:      "ALI_ACCESS_KEYID",
			Usage:       "The access key of accessing aliyun API.",
			Destination: &Aliyun.AccessKey,
		},
		cli.StringFlag{
			Name:        "ali-access-secret",
			EnvVar:      "ALI_ACCESS_SECRET",
			Usage:       "The access secret of accessing aliyun API",
			Destination: &Aliyun.AccessSecret,
		},
		cli.StringFlag{
			Name:        "ali-sign-name",
			EnvVar:      "ALI_SIGN_NAME",
			Usage:       "The signature of sms",
			Destination: &Aliyun.SMSSignName,
		},
		cli.StringFlag{
			Name:        "ali-template-code",
			EnvVar:      "ALI_TEMPLATE_CODE",
			Usage:       "The template code of aliyun sms",
			Destination: &Aliyun.SMSTemplateCode,
		},
	}
}
