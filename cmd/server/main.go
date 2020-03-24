package main

import (
	"flag"
	log "github.com/sirupsen/logrus"
	"gitlab.com/pliego/pxe-injector/pkg/tftp"
	"time"
)

func main() {
	flag.String("u", "root", "Specify username. Default is root")
	flag.String("p", "password", "Specify pass. Default is password")
	flag.Parse()

	log.Info("Hello World")

	s := new(tftp.Server)
	s.Start(tftp.ServerConfig{
		Address: ":69",
		Timeout: 5 * time.Second,
	})
}
