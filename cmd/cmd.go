package cmd

import (
	"CBCTF/internal/config"
	"CBCTF/internal/k8s"
	"CBCTF/internal/log"
	"flag"
	"os"
)

func init() {
	config.Init()
	log.Init()
}

func Cmd() {
	if len(os.Args) < 3 {
		run()
		return
	}
	cmd := flag.NewFlagSet("k8s", flag.ExitOnError)
	err := cmd.Parse(os.Args[2:])
	if err != nil {
		log.Logger.Fatalf("Failed to parse command: %s", err)
	}

	k8s.Init()
	switch os.Args[2] {
	case "init":
		k8s.InitResources()
	default:
		run()
	}
}
