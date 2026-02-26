package cmd

import (
	"CBCTF/internal/config"
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
)

func init() {
	i18n.Init()
	config.Init()
	log.Init()
}

func Cmd() {
	run()
}
