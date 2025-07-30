package config

import (
	_ "embed"
	"errors"
	"github.com/spf13/viper"
	"github.com/vmihailenco/msgpack/v5"
	"log"
	"os"
	"slices"
	"strings"
	"time"
)

var (
	allowedLogLevel     = []string{"DEBUG", "ERROR", "WARNING", "INFO"}
	allowedGinMode      = []string{"debug", "release", "test"}
	allowedGormLogLevel = []string{"INFO", "WARNING", "ERROR", "SILENT"}
)

type Config struct {
	Frontend  string `mapstructure:"frontend" json:"frontend"`     // 前端地址
	Backend   string `mapstructure:"backend" json:"backend"`       // 后端地址
	StaticURI string `mapstructure:"static_uri" json:"static_uri"` // 前端静态资源路径
	Path      string `mapstructure:"path" json:"path"`             // 数据存储路径

	Log struct {
		Level string `mapstructure:"level" json:"level"` // 日志级别：DEBUG, INFO, WARNING, ERROR
		Save  bool   `mapstructure:"save" json:"save"`   // 是否保存日志到文件
	} `mapstructure:"log" json:"log"`

	Gin struct {
		Mode   string `mapstructure:"mode" json:"mode"` // Gin 模式：debug, release, test
		Host   string `mapstructure:"host" json:"host"` // Gin 服务监听地址
		Port   int    `mapstructure:"port" json:"port"` // Gin 服务监听端口
		Upload struct {
			Max int `mapstructure:"max" json:"max"` // 上传文件最大大小（单位：MB）
		} `mapstructure:"upload" json:"upload"`
		Proxies   []string `mapstructure:"proxies" json:"proxies"` // 信任的代理服务器
		RateLimit struct {
			Whitelist []string `mapstructure:"whitelist" json:"whitelist"` // IP 白名单，不限制频率
		} `mapstructure:"rate" json:"rate"`
		Log struct {
			Whitelist []string `mapstructure:"whitelist" json:"whitelist"` // 日志白名单路径
		} `mapstructure:"log" json:"log"`
	} `mapstructure:"gin" json:"gin"`

	Gorm struct {
		MySQL struct {
			Host         string `mapstructure:"host" json:"host"`     // 数据库地址
			Port         int    `mapstructure:"port" json:"port"`     // 数据库端口
			User         string `mapstructure:"user" json:"user"`     // 数据库用户名
			Pwd          string `mapstructure:"pwd" json:"-"`         // 数据库密码
			DB           string `mapstructure:"db" json:"db"`         // 数据库名称
			MaxOpenConns int    `mapstructure:"mxopen" json:"mxopen"` // 最大连接数
			MaxIdleConns int    `mapstructure:"mxidle" json:"mxidle"` // 最大空闲连接数
		} `mapstructure:"mysql" json:"mysql"`
		Log struct {
			Level string `mapstructure:"level" json:"level"` // GORM 日志级别：INFO, WARNING, ERROR, SILENT
		} `mapstructure:"log" json:"log"`
	} `mapstructure:"gorm" json:"gorm"`

	Redis struct {
		Host string `mapstructure:"host" json:"host"` // Redis 地址
		Port int    `mapstructure:"port" json:"port"` // Redis 端口
		Pwd  string `mapstructure:"pwd" json:"-"`     // Redis 密码
	} `mapstructure:"redis" json:"redis"`

	K8S struct {
		Config          string `mapstructure:"config" json:"config"`
		Namespace       string `mapstructure:"namespace" json:"namespace"` // Kubernetes 命名空间
		ExternalNetwork struct {
			CIDR       string   `mapstructure:"cidr" json:"cidr"` // 外部网络 CIDR
			Gateway    string   `mapstructure:"gateway" json:"gateway"`
			Interface  string   `mapstructure:"interface" json:"interface"`
			ExcludeIPs []string `mapstructure:"exclude_ips" json:"exclude_ips"`
		} `mapstructure:"external_network" json:"external_network"`
		TCPDumpImage string `mapstructure:"tcpdump" json:"tcpdump"` // TCPDump 镜像
		Frpc         struct {
			On    bool   `mapstructure:"on" json:"on"`
			Image string `mapstructure:"image" json:"image"` // Frpc 镜像
			Frps  []struct {
				Host         string `mapstructure:"host" json:"host"`   // Frps 服务器地址
				Port         int    `mapstructure:"port" json:"port"`   // Frps 服务器端口
				Token        string `mapstructure:"token" json:"token"` // Frps 服务器 Token
				AllowedPorts []struct {
					From    int32   `mapstructure:"from" json:"from"`       // Frps 服务器允许的端口范围
					To      int32   `mapstructure:"to" json:"to"`           // Frps 服务器允许的端口范围
					Exclude []int32 `mapstructure:"exclude" json:"exclude"` // Frps 服务器排除的端口
				} `mapstructure:"allowed_ports" json:"allowed_ports"` // Frps 服务器允许的端口范围
			} `mapstructure:"frps" json:"frps"` // Frps 服务器列表
		} `mapstructure:"frpc" json:"frpc"`
		Nodes []string `mapstructure:"nodes" json:"nodes"` // Kubernetes 节点列表
	} `mapstructure:"k8s" json:"k8s"`

	NFS struct {
		Server  string `mapstructure:"server" json:"server"`
		Path    string `mapstructure:"path" json:"path"`
		Storage string `mapstructure:"storage" json:"storage"`
	} `mapstructure:"nfs" json:"nfs"`

	Email struct {
		Senders []struct {
			Addr string `mapstructure:"addr" json:"addr"` // 发件人地址
			Host string `mapstructure:"host" json:"host"` // SMTP 服务器地址
			Port int    `mapstructure:"port" json:"port"` // SMTP 服务器端口
			Pwd  string `mapstructure:"pwd" json:"-"`     // SMTP 服务器密码
		} `mapstructure:"senders" json:"senders"` // 发件人列表
	} `mapstructure:"email" json:"email"`
}

var Env *Config
var last time.Time

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
	if !slices.Contains(allowedLogLevel, strings.ToUpper(Env.Log.Level)) {
		Env.Log.Level = "INFO"
	}
	if !slices.Contains(allowedGinMode, strings.ToLower(Env.Gin.Mode)) {
		Env.Gin.Mode = "release"
	}
	if !slices.Contains(allowedGormLogLevel, strings.ToUpper(Env.Gin.Mode)) {
		Env.Gin.Mode = "SILENT"
	}
	Env.Backend = strings.TrimSuffix(Env.Backend, "/")
	Env.Frontend = strings.TrimSuffix(Env.Frontend, "/")
}

// Save 保存配置, 用于动态刷新配置
func Save(env *Config) error {
	env.Backend = strings.TrimSuffix(env.Backend, "/")
	env.Frontend = strings.TrimSuffix(env.Frontend, "/")
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
