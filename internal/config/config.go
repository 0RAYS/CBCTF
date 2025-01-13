package config

import (
	"bytes"
	_ "embed"
	"errors"
	"github.com/spf13/viper"
	"log"
	"os"
)

var Env *viper.Viper

//go:embed default.yml
var defaultConf []byte

func Init() {
	Env = viper.New()
	Env.SetConfigType("yaml")
	Env.SetConfigName("config")
	Env.AddConfigPath(".")
	if err := Env.ReadInConfig(); err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if errors.As(err, &configFileNotFoundError) {
			if err := os.WriteFile("./config.yml", defaultConf, 0666); err != nil {
				log.Panicf("Failed to init config: %s", err)
			}
			_ = Env.ReadConfig(bytes.NewReader(defaultConf))
		}
	}
	check()
}

func check() {
	for _, v := range Env.AllKeys() {
		if Env.Get(v) == "" || Env.Get(v) == 0 {
			log.Panicf("Invalid environment variable %s: %v", v, Env.Get(v))
		}
	}
}
