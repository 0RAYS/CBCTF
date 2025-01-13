package main

import (
	"CBCTF/router"
	"github.com/sirupsen/logrus"
	"log"
)

func init() {

}

func main() {
	logrus.Info("test")
	if err := router.Init().Run(); err != nil {
		log.Fatalf("Failed to run router: %s", err)
	}
}
