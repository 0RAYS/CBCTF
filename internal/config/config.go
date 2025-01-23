package config

import (
	"bytes"
	_ "embed"
	"errors"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"log"
	"os"
)

type Config struct {
	Log struct {
		Level string `mapstructure:"level" json:"level"` // 日志级别：DEBUG, INFO, WARNING, ERROR
		Save  bool   `mapstructure:"save" json:"save"`   // 是否保存日志到文件
	} `mapstructure:"log" json:"log"`

	Gin struct {
		Mode   string `mapstructure:"mode" json:"mode"` // Gin 模式：debug, release, test
		Host   string `mapstructure:"host" json:"host"` // Gin 服务监听地址
		Port   int    `mapstructure:"port" json:"port"` // Gin 服务监听端口
		Upload struct {
			Path string `mapstructure:"path" json:"path"` // 上传文件路径
			Max  int    `mapstructure:"max" json:"max"`   // 上传文件最大大小（单位：MB）
		} `mapstructure:"upload" json:"upload"`
	} `mapstructure:"gin" json:"gin"`

	Gorm struct {
		File string `mapstructure:"file" json:"file"` // 数据库文件路径
		Log  struct {
			Level string `mapstructure:"level" json:"level"` // GORM 日志级别：INFO, WARNING, ERROR, SILENT
		} `mapstructure:"log" json:"log"`
	} `mapstructure:"gorm" json:"gorm"`

	Redis struct {
		Addr    string `mapstructure:"addr" json:"addr"`       // Redis 地址
		Pwd     string `mapstructure:"pwd" json:"pwd"`         // Redis 密码
		Timeout uint   `mapstructure:"timeout" json:"timeout"` // Redis 连接超时时间（单位：毫秒）
	} `mapstructure:"redis" json:"redis"`

	Url string `mapstructure:"url" json:"url"` // 主机地址
}

var Env *Config

//go:embed default.yml
var defaultConf []byte

func Init() {
	viper.SetConfigType("yaml")
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if errors.As(err, &configFileNotFoundError) {
			if err := os.WriteFile("./config.yml", defaultConf, 0666); err != nil {
				log.Panicf("Failed to init config: %s", err)
			}
			_ = viper.ReadConfig(bytes.NewReader(defaultConf))
		}
	}
	if err := viper.Unmarshal(&Env); err != nil {
		log.Panicf("error unmarshalling config: %s", err)
	}
}

func Watch(onChange func()) {
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Printf("Config file changed: %s", e.Name)
		onChange()
	})
}
