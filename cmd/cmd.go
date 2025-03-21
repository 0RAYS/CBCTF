package cmd

import (
	"CBCTF/internal/config"
	"CBCTF/internal/k8s"
	"CBCTF/internal/log"
	"flag"
)

func init() {
	config.Init()
	log.Init()
}

func Cmd() {
	cmd := flag.NewFlagSet("k8s", flag.ExitOnError)
	init := cmd.Bool("init", false, "init k8s")
	check := cmd.Bool("check", false, "check k8s")

	flag.Parse()
	args := flag.Args()

	switch len(args) {
	case 0, 1:
		run()
	case 2:
		err := cmd.Parse(args[1:])
		if err != nil {
			log.Logger.Fatalf("Failed to parse command: %v", err)
		}
		k8s.Init(false)
		if *init {
			k8s.InitResources()
		} else if *check {
			k8s.CheckPermission()
		} else {
			run()
		}

	}
}
