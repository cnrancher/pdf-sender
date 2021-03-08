package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/cnrancher/pdf-sender/pkg/apis"
	"github.com/cnrancher/pdf-sender/pkg/limiter"
	"github.com/cnrancher/pdf-sender/pkg/types"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"

	_ "github.com/go-sql-driver/mysql"
)

var (
	VERSION = "v0.0.0-dev"
	port    int
)

func main() {
	app := cli.NewApp()
	app.Name = "pdf-sender"
	app.Version = VERSION
	app.Usage = "Send pdf documents to our lovely users."
	app.Flags = []cli.Flag{
		cli.IntFlag{
			Name:        "port,p",
			EnvVar:      "HTTP_PORT",
			Value:       8080,
			Destination: &port,
		},
		cli.BoolFlag{
			Name: "debug",
		},
	}
	app.Flags = append(app.Flags, apis.GetFlags()...)
	app.Flags = append(app.Flags, types.GetFlags()...)
	app.Action = func(ctx *cli.Context) error {
		logrus.Infof("server running, listening at :%d\n", port)
		return http.ListenAndServe(fmt.Sprintf(":%d", port), limiter.IPLimitMiddleware(apis.RegisterAPIs().ServeMux))
	}
	app.Before = before

	if err := app.Run(os.Args); err != nil {
		logrus.Fatal(err)
	}
}

func before(ctx *cli.Context) error {
	if ctx.Bool("debug") {
		logrus.SetLevel(logrus.DebugLevel)
	}

	if err := types.InitAliyunClients(ctx); err != nil {
		return err
	}

	if err := apis.InitEmailBody(ctx); err != nil {
		return err
	}

	if err := apis.ConnectMysql(); err != nil {
		return err
	}

	apis.StartCorn(ctx)

	return nil
}
