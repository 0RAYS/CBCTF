package config

import (
	"bytes"
	_ "embed"
	"encoding/json"
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
		Type   string `mapstructure:"type" json:"type"` // 数据库类型：sqlite, mysql
		SQLite struct {
			File string `mapstructure:"file" json:"file"` // 数据库文件路径
		} `mapstructure:"sqlite" json:"sqlite"`
		MySQL struct {
			Host string `mapstructure:"host" json:"host"` // 数据库地址
			Port int    `mapstructure:"port" json:"port"` // 数据库端口
			User string `mapstructure:"user" json:"user"` //
			Pwd  string `mapstructure:"pwd" json:"pwd"`   // 数据库密码
			DB   string `mapstructure:"db" json:"db"`     // 数据库名称
		} `mapstructure:"mysql" json:"mysql"`
		Log struct {
			Level string `mapstructure:"level" json:"level"` // GORM 日志级别：INFO, WARNING, ERROR, SILENT
		} `mapstructure:"log" json:"log"`
	} `mapstructure:"gorm" json:"gorm"`

	Redis struct {
		Addr    string `mapstructure:"addr" json:"addr"`       // Redis 地址
		Pwd     string `mapstructure:"pwd" json:"pwd"`         // Redis 密码
		Timeout uint   `mapstructure:"timeout" json:"timeout"` // Redis 连接超时时间（单位：毫秒）
	} `mapstructure:"redis" json:"redis"`

	Frontend string `mapstructure:"frontend" json:"frontend"` // 前端地址
	Backend  string `mapstructure:"backend" json:"backend"`   // 后端地址
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

func Save(env Config) error {
	config := make(map[string]interface{})
	data, err := json.Marshal(env)
	if err != nil {
		log.Panicf("Failed to marshal Env to JSON: %s", err)
		return err
	}
	if err = json.Unmarshal(data, &config); err != nil {
		log.Panicf("Failed to unmarshal JSON to map: %s", err)
		return err
	}
	if err := viper.MergeConfigMap(config); err != nil {
		log.Panicf("Failed to merge Env to viper: %s", err)
		return err
	}
	if err := viper.WriteConfig(); err != nil {
		log.Panicf("Failed to save config: %s", err)
		return err
	}
	return nil
}

func Watch(onChange func()) {
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Printf("Config file changed: %s", e.Name)
		onChange()
	})
}
