package main

import (
	"flag"

	"github.com/cnrancher/pdf-sender/pkg/server"
	"github.com/sirupsen/logrus"
)

var (
	port = flag.Int("port", 8080, "listen port")
)

func main() {
	flag.Parse()
	if err := server.New(*port).Run(); err != nil {
		logrus.Fatalf("server run fatal:%v", err)
	}
}
