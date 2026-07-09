package config

import (
	"CBCTF/internal/log"
	"CBCTF/internal/utils"
	"bytes"
	_ "embed"
	"errors"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

type FrpsConfig struct {
	Host    string `json:"host" binding:"required,hostname|ip"`
	Token   string `json:"token"`
	Allowed []struct {
		From    int32   `json:"from" binding:"required,gte=0,lte=65535"`
		To      int32   `json:"to" binding:"required,gte=0,lte=65535"`
		Exclude []int32 `json:"exclude" binding:"dive,gte=0,lte=65535"`
	} `json:"allowed" binding:"required,dive"`
	Port int32 `json:"port" binding:"required,gte=0,lte=65535"`
}

type Config struct {
	Gorm struct {
		Postgres struct {
			Host         string `mapstructure:"host" json:"host"`
			Port         uint   `mapstructure:"port" json:"port"`
			User         string `mapstructure:"user" json:"user"`
			Pwd          string `mapstructure:"pwd" json:"pwd"`
			DB           string `mapstructure:"db" json:"db"`
			SSLMode      bool   `mapstructure:"sslmode" json:"sslmode"`
			MaxOpenConns int    `mapstructure:"mxopen" json:"mxopen"`
			MaxIdleConns int    `mapstructure:"mxidle" json:"mxidle"`
		} `mapstructure:"postgres" json:"postgres"`
		Log struct {
			Level string `mapstructure:"level" json:"level"`
		} `mapstructure:"log" json:"log"`
	} `mapstructure:"gorm" json:"gorm"`
	Redis struct {
		Host string `mapstructure:"host" json:"host"`
		Port uint   `mapstructure:"port" json:"port"`
		Pwd  string `mapstructure:"pwd" json:"pwd"`
	} `mapstructure:"redis" json:"redis"`
	Host string `mapstructure:"host" json:"host"`
	Path string `mapstructure:"path" json:"path"`
	Gin  struct {
		Mode   string `mapstructure:"mode" json:"mode"`
		Host   string `mapstructure:"host" json:"host"`
		Port   uint   `mapstructure:"port" json:"port"`
		Upload struct {
			Picture   int `mapstructure:"picture" json:"picture"`
			Challenge int `mapstructure:"challenge" json:"challenge"`
			Writeup   int `mapstructure:"writeup" json:"writeup"`
		} `mapstructure:"upload" json:"upload"`
		Proxies   []string `mapstructure:"proxies" json:"proxies"`
		RateLimit struct {
			Global    int      `mapstructure:"global" json:"global"`
			Whitelist []string `mapstructure:"whitelist" json:"whitelist"`
		} `mapstructure:"ratelimit" json:"ratelimit"`
		Origins []string `mapstructure:"origins" json:"origins"`
		Log     struct {
			Whitelist []string `mapstructure:"whitelist" json:"whitelist"`
		} `mapstructure:"log" json:"log"`
		JWT struct {
			Secret string `mapstructure:"secret" json:"secret"`
		} `mapstructure:"jwt" json:"jwt"`
		PProf struct {
			Whitelist []string `mapstructure:"whitelist" json:"whitelist"`
		} `mapstructure:"pprof" json:"pprof"`
		Metrics struct {
			Whitelist []string `mapstructure:"whitelist" json:"whitelist"`
		} `mapstructure:"metrics" json:"metrics"`
	} `mapstructure:"gin" json:"gin"`
	K8S struct {
		Namespace    string `mapstructure:"namespace" json:"namespace"`
		CaptureImage string `mapstructure:"capture" json:"capture"`
		Frp          struct {
			On         bool         `mapstructure:"on" json:"on"`
			FrpcImage  string       `mapstructure:"frpc" json:"frpc"`
			NginxImage string       `mapstructure:"nginx" json:"nginx"`
			Frps       []FrpsConfig `mapstructure:"frps" json:"frps"`
		} `mapstructure:"frp" json:"frp"`
	} `mapstructure:"k8s" json:"k8s"`
	Cheat struct {
		IP struct {
			Whitelist []string `mapstructure:"whitelist" json:"whitelist"`
		} `mapstructure:"ip" json:"ip"`
	} `mapstructure:"cheat" json:"cheat"`
	Webhook struct {
		Whitelist []string `mapstructure:"whitelist" json:"whitelist"`
	} `mapstructure:"webhook" json:"webhook"`
	AsyncQ struct {
		Log struct {
			Level string `mapstructure:"level" json:"level"`
		} `mapstructure:"log" json:"log"`
		Queues struct {
			Victim     int `mapstructure:"victim" json:"victim"`
			Traffic    int `mapstructure:"traffic" json:"traffic"`
			Generator  int `mapstructure:"generator" json:"generator"`
			Attachment int `mapstructure:"attachment" json:"attachment"`
			Email      int `mapstructure:"email" json:"email"`
			Webhook    int `mapstructure:"webhook" json:"webhook"`
			Image      int `mapstructure:"image" json:"image"`
		} `mapstructure:"queues" json:"queues"`
	} `mapstructure:"asynq" json:"asynq"`
	Registration struct {
		Enabled      bool `mapstructure:"enabled" json:"enabled"`
		DefaultGroup uint `mapstructure:"default_group" json:"default_group"`
	} `mapstructure:"registration" json:"registration"`
}

var Env *Config

//go:embed default.yaml
var defaultConf []byte

func Init(path string) {
	v := viper.New()
	v.SetConfigType("yaml")
	if err := v.ReadConfig(bytes.NewReader(defaultConf)); err != nil {
		log.Logger.Fatalf("Failed to read default config: %s", err)
	}
	if path != "" {
		v.SetConfigFile(path)
		var notFound viper.ConfigFileNotFoundError
		if err := v.MergeInConfig(); err != nil && !errors.As(err, &notFound) {
			log.Logger.Fatalf("Failed to read config: %s", err)
		}
	}
	if err := v.Unmarshal(&Env); err != nil {
		log.Logger.Fatalf("error unmarshalling config: %s", err)
	}
	tidy()
}

func tidy() {
	Env.Host = strings.TrimSuffix(Env.Host, "/")
	Env.Path = filepath.Clean(Env.Path)
	if Env.Gin.JWT.Secret == "" {
		Env.Gin.JWT.Secret = utils.UUID()
	}
}
