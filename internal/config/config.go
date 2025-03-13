package config

import (
	_ "embed"
	"encoding/json"
	"errors"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"log"
	"os"
	"reflect"
	"strings"
	"time"
)

type Sender struct {
	Address  string `json:"address"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Password string `json:"password" secret:"true"`
}

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
		Proxies []string `mapstructure:"proxies" json:"proxies"` // 信任的代理服务器
		Magic   struct {
			Whitelist []string `mapstructure:"whitelist" json:"whitelist"`
		} `mapstructure:"magic" json:"magic"`
	} `mapstructure:"gin" json:"gin"`

	Gorm struct {
		MySQL struct {
			Host         string `mapstructure:"host" json:"host"`             // 数据库地址
			Port         int    `mapstructure:"port" json:"port"`             // 数据库端口
			User         string `mapstructure:"user" json:"user"`             //
			Pwd          string `mapstructure:"pwd" json:"pwd" secret:"true"` // 数据库密码
			DB           string `mapstructure:"db" json:"db"`                 // 数据库名称
			MaxOpenConns int    `mapstructure:"mxopen" json:"mxopen"`         // 最大连接数
			MaxIdleConns int    `mapstructure:"mxidle" json:"mxidle"`         // 最大空闲连接数
		} `mapstructure:"mysql" json:"mysql"`
		Log struct {
			Level string `mapstructure:"level" json:"level"` // GORM 日志级别：INFO, WARNING, ERROR, SILENT
		} `mapstructure:"log" json:"log"`
	} `mapstructure:"gorm" json:"gorm"`

	Redis struct {
		On      bool   `mapstructure:"on" json:"on"`                 // Redis 开关
		Addr    string `mapstructure:"addr" json:"addr"`             // Redis 地址
		Pwd     string `mapstructure:"pwd" json:"pwd" secret:"true"` // Redis 密码
		Timeout uint   `mapstructure:"timeout" json:"timeout"`       // Redis 连接超时时间（单位：毫秒）
	} `mapstructure:"redis" json:"redis"`

	K8S struct {
		Config       string   `mapstructure:"config" json:"config"`       // Kubernetes 配置文件路径
		Master       string   `mapstructure:"master" json:"master"`       // Kubernetes Master 地址
		Namespace    string   `mapstructure:"namespace" json:"namespace"` // Kubernetes 命名空间
		TCPDumpImage string   `mapstructure:"tcpdump" json:"tcpdump"`     // TCPDump 镜像
		Nodes        []string `mapstructure:"nodes" json:"nodes"`         // Kubernetes 节点列表
	} `mapstructure:"k8s" json:"k8s"`

	Email struct {
		Senders []Sender `mapstructure:"senders" json:"senders"` // 发件人列表
	} `mapstructure:"email" json:"email"`

	Frontend string `mapstructure:"frontend" json:"frontend"` // 前端地址
	Backend  string `mapstructure:"backend" json:"backend"`   // 后端地址
}

func MaskSecrets(input interface{}) interface{} {
	val := reflect.ValueOf(input)

	switch val.Kind() {
	case reflect.Ptr: // 指针类型
		if val.IsNil() {
			return nil
		}
		return MaskSecrets(val.Elem().Interface())

	case reflect.Struct: // 结构体类型
		result := make(map[string]interface{})
		typ := val.Type()
		for i := 0; i < val.NumField(); i++ {
			field := typ.Field(i)
			jsonKey := field.Tag.Get("json")
			if jsonKey == "" {
				jsonKey = field.Name
			}

			if field.Tag.Get("secret") == "true" {
				result[jsonKey] = "******"
			} else {
				result[jsonKey] = MaskSecrets(val.Field(i).Interface())
			}
		}
		return result

	case reflect.Slice: // 切片类型
		length := val.Len()
		sliceResult := make([]interface{}, length)
		for i := 0; i < length; i++ {
			sliceResult[i] = MaskSecrets(val.Index(i).Interface())
		}
		return sliceResult

	case reflect.Map: // 映射类型
		mapResult := make(map[string]interface{})
		for _, key := range val.MapKeys() {
			mapResult[key.String()] = MaskSecrets(val.MapIndex(key).Interface())
		}
		return mapResult

	default: // 其他类型（string、int、bool...）
		return input
	}
}

func (c Config) MarshalJSON() ([]byte, error) {
	type Alias Config
	tmp := Alias(c)
	data := MaskSecrets(tmp)
	return json.Marshal(data)
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
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if errors.As(err, &configFileNotFoundError) {
			if err := os.WriteFile("./config.yml", defaultConf, 0666); err != nil {
				log.Panicf("Failed to init config: %s", err)
			}
			log.Panicf("Please configure the config.yml file and restart the program")
		}
	}
	if err := viper.Unmarshal(&Env); err != nil {
		log.Panicf("error unmarshalling config: %s", err)
	}
	Env.Backend = strings.TrimSuffix(Env.Backend, "/")
	Env.Frontend = strings.TrimSuffix(Env.Frontend, "/")
}

// Save 保存配置，用于动态刷新配置
//func Save(env Config) error {
//	env.Backend = strings.TrimSuffix(env.Backend, "/")
//	env.Frontend = strings.TrimSuffix(env.Frontend, "/")
//	config := make(map[string]interface{})
//	data, err := json.Marshal(env)
//	if err != nil {
//		log.Panicf("Failed to marshal Env to JSON: %s", err)
//		return err
//	}
//	if err = json.Unmarshal(data, &config); err != nil {
//		log.Panicf("Failed to unmarshal JSON to map: %s", err)
//		return err
//	}
//	if err := viper.MergeConfigMap(config); err != nil {
//		log.Panicf("Failed to merge Env to viper: %s", err)
//		return err
//	}
//	if err := viper.WriteConfig(); err != nil {
//		log.Panicf("Failed to save config: %s", err)
//		return err
//	}
//	return nil
//}

// Watch 监听配置文件变化
func Watch(onChange func()) {
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		if time.Since(last) < time.Second {
			return
		}
		last = time.Now()
		log.Printf("Config file changed: %s", e.Name)
		onChange()
	})
}
