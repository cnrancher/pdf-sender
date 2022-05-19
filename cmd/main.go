package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/cnrancher/pdf-sender/pkg/apis"
	"github.com/cnrancher/pdf-sender/pkg/email"
	"github.com/cnrancher/pdf-sender/pkg/limiter"
	"github.com/cnrancher/pdf-sender/pkg/types"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"

	_ "github.com/go-sql-driver/mysql"
)

var (
	VERSION = "v0.0.0-dev"
)

func main() {
	app := cli.NewApp()
	app.Name = "pdf-sender"
	app.Version = VERSION
	app.Usage = "Send pdf documents to our lovely users."
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "debug",
			Usage: "Enable Debug log for pdf sender",
		},
		cli.BoolFlag{
			Name:  "dry-run",
			Usage: "dry-run mode will enable debug log and won't send anything",
		},
		cli.StringFlag{
			Name:     "config-file,f",
			Required: true,
			EnvVar:   "CONFIG_FILE",
			Value:    "/etc/pdf-sender.yml",
		},
	}
	app.Commands = append(app.Commands, types.GetConfigCommand())
	app.Commands = append(app.Commands, getRunCommand())

	if err := app.Run(os.Args); err != nil {
		logrus.Fatal(err)
	}
}

func before(ctx *cli.Context) error {
	if err := types.MergeConfig(ctx); err != nil {
		return err
	}

	types.SetRunStatus(ctx.GlobalBool("debug"), ctx.GlobalBool("dry-run"))

	if err := types.ConnectDB(); err != nil {
		return err
	}

	if err := types.InitAliyunClients(ctx); err != nil {
		return err
	}

	if err := email.InitEmailBody(ctx); err != nil {
		return err
	}

	if err := types.Validate(); err != nil {
		return err
	}

	return apis.StartCorn(ctx)
}

func getRunCommand() cli.Command {
	cmd := cli.Command{
		Name:  "run",
		Usage: "Run pdf sender server",
		Flags: types.GetFlags(),
		Action: func(ctx *cli.Context) error {
			logrus.Infof("server running, listening at :%d\n", types.Config.Port)
			return http.ListenAndServe(fmt.Sprintf(":%d", types.Config.Port), limiter.IPLimitMiddleware(apis.RegisterAPIs().ServeMux))
		},
		Before: before,
	}
	return cmd
}
