package cmd

import (
	"CBCTF/internel/config"
	//"CBCTF/internel/k8s"
	"CBCTF/internel/log"
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
	}
	cmd := flag.NewFlagSet("k8s", flag.ExitOnError)
	err := cmd.Parse(os.Args[2:])
	if err != nil {
		log.Logger.Fatalf("Failed to parse command: %v", err)
	}

	//k8s.Init(false)
	switch os.Args[2] {
	//case "init":
	//	k8s.InitResources()
	//case "check":
	//	k8s.CheckPermission()
	default:
		run()
	}
}
