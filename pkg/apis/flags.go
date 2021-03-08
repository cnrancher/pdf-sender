package apis

import (
	"github.com/urfave/cli"
)

func GetFlags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:        "pdf-file-url-prefix",
			EnvVar:      "PDF_FILE_URL_PREFIX",
			Usage:       "The url prefix for all pdf files if the files hosted in a HTTP server.",
			Destination: &fileURLPrefix,
		},
		cli.StringSliceFlag{
			Name:   "published-docs",
			EnvVar: "PUBLISHED_DOCS",
			Value:  &cli.StringSlice{"rancher2", "rke", "k3s", "octopus", "harvester"},
		},

		cli.StringFlag{
			Name:        "smtp-user",
			EnvVar:      "SMTP_USER",
			Usage:       "The smtp server user for sending emails.",
			Destination: &smtpUser,
		},
		cli.StringFlag{
			Name:        "smtp-pwd",
			EnvVar:      "SMTP_PWD",
			Usage:       "the smtp user password for sending emails.",
			Destination: &smtpPWD,
		},
		cli.StringFlag{
			Name:        "smtp-endpoint",
			EnvVar:      "SMTP_ENDPOINT",
			Usage:       "The smtp server address or domain name.",
			Destination: &smtpEndpoint,
		},
		cli.IntFlag{
			Name:        "smtp-port",
			EnvVar:      "SMTP_PORT",
			Usage:       "The smtp server port",
			Destination: &smtpPort,
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
			Destination: &DayCronJob,
			Value:       "30 1 * * ?",
		},
		cli.StringFlag{
			Name:        "stmp-monthly-cron,m",
			EnvVar:      "MON_CRON",
			Usage:       "The monthly cron string for sending statistics.",
			Destination: &MonCronJob,
			Value:       "0 2 1 * ?",
		},
		cli.StringFlag{
			Name:        "smtp-sender-email,s",
			Usage:       "The sender email address.",
			EnvVar:      "SENDER_EMAIL",
			Destination: &senderEmail,
		},

		cli.StringFlag{
			Name:        "db-host-ip",
			EnvVar:      "DB_HOST_IP",
			Usage:       "The backend mysql hostname or IP.",
			Required:    true,
			Destination: &dbhostip,
		},
		cli.IntFlag{
			Name:        "db-port",
			EnvVar:      "DB_PORT",
			Usage:       "The backend mysql port.",
			Destination: &dbport,
			Value:       3306,
		},
		cli.StringFlag{
			Name:        "db-name",
			EnvVar:      "DB_NAME",
			Usage:       "The db name of mysql backend.",
			Destination: &dbname,
			Value:       "pdf",
		},
		cli.StringFlag{
			Name:        "db-username",
			EnvVar:      "DB_USERNAME",
			Usage:       "The username of mysql backend.",
			Required:    true,
			Destination: &dbusername,
		},
		cli.StringFlag{
			Name:        "db-password",
			EnvVar:      "DB_PASSWORD",
			Usage:       "The password of mysql user.",
			Required:    true,
			Destination: &dbpassword,
		},
	}
}
