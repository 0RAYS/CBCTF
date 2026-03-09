package cmd

import (
	"CBCTF/internal/i18n"
	"flag"
	"fmt"
	"os"
)

var configPath string

func init() {
	i18n.Init()
}

func Cmd() {
	fs := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	fs.StringVar(&configPath, "c", "config.yaml", "Path to config file")
	fs.Usage = func() {
		_, _ = fmt.Fprintf(os.Stderr, "Usage: %s [options] [command]\n\n", os.Args[0])
		_, _ = fmt.Fprintln(os.Stderr, "Options:")
		fs.PrintDefaults()
	}
	_ = fs.Parse(os.Args[1:])

	preInit()
	run()
}
