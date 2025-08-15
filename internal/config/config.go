package config

import (
	_ "embed"
	"errors"
	"log"
	"os"
	"slices"
	"strings"

	"github.com/spf13/viper"
	"github.com/vmihailenco/msgpack/v5"
)

var (
	allowedLogLevel       = []string{"DEBUG", "ERROR", "WARNING", "INFO"}
	allowedAsyncQLogLevel = []string{"DEBUG", "ERROR", "WARNING", "INFO"}
	allowedGinMode        = []string{"debug", "release", "test"}
	allowedGormLogLevel   = []string{"INFO", "WARNING", "ERROR", "SILENT"}
)

type Config struct {
	Backend  string `mapstructure:"backend" json:"backend" msgpack:"backend"`    // 后端地址
	Frontend string `mapstructure:"frontend" json:"frontend" msgpack:"frontend"` // 前端地址
	Path     string `mapstructure:"path" json:"path" msgpack:"path"`             // 数据存储路径

	Log struct {
		Level string `mapstructure:"level" json:"level" msgpack:"level"` // 日志级别：DEBUG, INFO, WARNING, ERROR
		Save  bool   `mapstructure:"save" json:"save" msgpack:"save"`    // 是否保存日志到文件
	} `mapstructure:"log" json:"log" msgpack:"log"`

	AsyncQ struct {
		Level       string `mapstructure:"level" json:"level" msgpack:"level"`
		Concurrency int    `mapstructure:"concurrency" json:"concurrency" msgpack:"concurrency"`
	} `mapstructure:"asynq" json:"asynq" msgpack:"asynq"`

	Gin struct {
		Mode   string `mapstructure:"mode" json:"mode" msgpack:"mode"` // Gin 模式：debug, release, test
		Host   string `mapstructure:"host" json:"host" msgpack:"host"` // Gin 服务监听地址
		Port   int    `mapstructure:"port" json:"port" msgpack:"port"` // Gin 服务监听端口
		Upload struct {
			Max int `mapstructure:"max" json:"max" msgpack:"max"` // 上传文件最大大小（单位：MB）
		} `mapstructure:"upload" json:"upload" msgpack:"upload"`
		Proxies   []string `mapstructure:"proxies" json:"proxies" msgpack:"proxies"` // 信任的代理服务器
		RateLimit struct {
			Whitelist []string `mapstructure:"whitelist" json:"whitelist" msgpack:"whitelist"` // IP 白名单，不限制频率
		} `mapstructure:"rate" json:"rate" msgpack:"rate"`
		Log struct {
			Whitelist []string `mapstructure:"whitelist" json:"whitelist" msgpack:"whitelist"` // 日志白名单路径
		} `mapstructure:"log" json:"log" msgpack:"log"`
	} `mapstructure:"gin" json:"gin" msgpack:"gin"`

	Gorm struct {
		MySQL struct {
			Host         string `mapstructure:"host" json:"host" msgpack:"host"`       // 数据库地址
			Port         int    `mapstructure:"port" json:"port" msgpack:"port"`       // 数据库端口
			User         string `mapstructure:"user" json:"user" msgpack:"user"`       // 数据库用户名
			Pwd          string `mapstructure:"pwd" json:"-" msgpack:"pwd"`            // 数据库密码
			DB           string `mapstructure:"db" json:"db" msgpack:"db"`             // 数据库名称
			MaxOpenConns int    `mapstructure:"mxopen" json:"mxopen" msgpack:"mxopen"` // 最大连接数
			MaxIdleConns int    `mapstructure:"mxidle" json:"mxidle" msgpack:"mxidle"` // 最大空闲连接数
		} `mapstructure:"mysql" json:"mysql" msgpack:"mysql"` // MySQL 数据库配置
		Log struct {
			Level string `mapstructure:"level" json:"level" msgpack:"level"` // GORM 日志级别：INFO, WARNING, ERROR, SILENT
		} `mapstructure:"log" json:"log" msgpack:"log"`
	} `mapstructure:"gorm" json:"gorm" msgpack:"gorm"`

	Redis struct {
		Host string `mapstructure:"host" json:"host" msgpack:"host"` // Redis 地址
		Port int    `mapstructure:"port" json:"port" msgpack:"port"` // Redis 端口
		Pwd  string `mapstructure:"pwd" json:"-" msgpack:"pwd"`      // Redis 密码
	} `mapstructure:"redis" json:"redis" msgpack:"redis"`

	K8S struct {
		Config          string `mapstructure:"config" json:"config" msgpack:"config"`
		Namespace       string `mapstructure:"namespace" json:"namespace" msgpack:"namespace"` // Kubernetes 命名空间
		ExternalNetwork struct {
			CIDR       string   `mapstructure:"cidr" json:"cidr" msgpack:"cidr"` // 外部网络 CIDR
			Gateway    string   `mapstructure:"gateway" json:"gateway" msgpack:"gateway"`
			Interface  string   `mapstructure:"interface" json:"interface" msgpack:"interface"`
			ExcludeIPs []string `mapstructure:"exclude_ips" json:"exclude_ips" msgpack:"exclude_ips"`
		} `mapstructure:"external_network" json:"external_network" msgpack:"external_network"`
		TCPDumpImage string `mapstructure:"tcpdump" json:"tcpdump" msgpack:"tcpdump"` // TCPDump 镜像
		Frpc         struct {
			On    bool   `mapstructure:"on" json:"on" msgpack:"on"`
			Image string `mapstructure:"image" json:"image" msgpack:"image"` // Frpc 镜像
			Frps  []struct {
				Host         string `mapstructure:"host" json:"host" msgpack:"host"`    // Frps 服务器地址
				Port         int    `mapstructure:"port" json:"port" msgpack:"port"`    // Frps 服务器端口
				Token        string `mapstructure:"token" json:"token" msgpack:"token"` // Frps 服务器 Token
				AllowedPorts []struct {
					From    int32   `mapstructure:"from" json:"from" msgpack:"from"`          // Frps 服务器允许的端口范围
					To      int32   `mapstructure:"to" json:"to" msgpack:"to"`                // Frps 服务器允许的端口范围
					Exclude []int32 `mapstructure:"exclude" json:"exclude" msgpack:"exclude"` // Frps 服务器排除的端口
				} `mapstructure:"allowed_ports" json:"allowed_ports" msgpack:"allowed_ports"` // Frps 服务器允许的端口范围
			} `mapstructure:"frps" json:"frps" msgpack:"frps"` // Frps 服务器列表
		} `mapstructure:"frpc" json:"frpc" msgpack:"frpc"`
		Nodes           []string `mapstructure:"nodes" json:"nodes" msgpack:"nodes"` // Kubernetes 节点列表
		GeneratorWorker int      `mapstructure:"generator_worker" json:"generator_worker" msgpack:"generator_worker"`
	} `mapstructure:"k8s" json:"k8s" msgpack:"k8s"`

	NFS struct {
		Server  string `mapstructure:"server" json:"server" msgpack:"server"`
		Path    string `mapstructure:"path" json:"path" msgpack:"path"`
		Storage string `mapstructure:"storage" json:"storage" msgpack:"storage"`
	} `mapstructure:"nfs" json:"nfs" msgpack:"nfs"`

	Email struct {
		Senders []struct {
			Addr string `mapstructure:"addr" json:"addr" msgpack:"addr"` // 发件人地址
			Host string `mapstructure:"host" json:"host" msgpack:"host"` // SMTP 服务器地址
			Port int    `mapstructure:"port" json:"port" msgpack:"port"` // SMTP 服务器端口
			Pwd  string `mapstructure:"pwd" json:"-" msgpack:"pwd"`      // SMTP 服务器密码
		} `mapstructure:"senders" json:"senders" msgpack:"senders"` // 发件人列表
	} `mapstructure:"email" json:"email" msgpack:"email"`
}

var Env *Config

//go:embed default.yml
var defaultConf []byte

// Init 初始化配置
func Init() {
	viper.SetConfigType("yaml")
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		if errors.As(err, &viper.ConfigFileNotFoundError{}) {
			if err := os.WriteFile("./config.yml", defaultConf, 0666); err != nil {
				log.Panicf("Failed to init config: %s", err)
			}
			log.Fatalf("Please configure the config.yml file and restart the program")
		}
	}
	if err := viper.Unmarshal(&Env); err != nil {
		log.Panicf("error unmarshalling config: %s", err)
	}
	tidy()
}

func tidy() {
	Env.Log.Level = strings.ToUpper(Env.Log.Level)
	if !slices.Contains(allowedLogLevel, Env.Log.Level) {
		Env.Log.Level = "INFO"
	}
	Env.AsyncQ.Level = strings.ToUpper(Env.AsyncQ.Level)
	if !slices.Contains(allowedAsyncQLogLevel, Env.AsyncQ.Level) {
		Env.AsyncQ.Level = "WARNING"
	}
	Env.Gin.Mode = strings.ToLower(Env.Gin.Mode)
	if !slices.Contains(allowedGinMode, Env.Gin.Mode) {
		Env.Gin.Mode = "release"
	}
	Env.Gorm.Log.Level = strings.ToUpper(Env.Gorm.Log.Level)
	if !slices.Contains(allowedGormLogLevel, Env.Gorm.Log.Level) {
		Env.Gorm.Log.Level = "SILENT"
	}
	Env.Backend = strings.TrimSuffix(Env.Backend, "/")
	Env.Frontend = strings.TrimSuffix(Env.Frontend, "/")
}

// Save 保存配置, 用于动态刷新配置
func Save(env *Config) error {
	config := make(map[string]any)
	data, err := msgpack.Marshal(env)
	if err != nil {
		log.Panicf("Failed to marshal Env to msgpack: %s", err)
		return err
	}
	if err = msgpack.Unmarshal(data, &config); err != nil {
		log.Panicf("Failed to unmarshal msgpack to map: %s", err)
		return err
	}
	if err = viper.MergeConfigMap(config); err != nil {
		log.Panicf("Failed to merge Env to viper: %s", err)
		return err
	}
	if err = viper.WriteConfig(); err != nil {
		log.Panicf("Failed to save config: %s", err)
		return err
	}
	return nil
}
