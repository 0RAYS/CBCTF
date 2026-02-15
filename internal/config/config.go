package config

import (
	_ "embed"
	"errors"
	"log"
	"os"
	"reflect"
	"strings"

	"github.com/mitchellh/mapstructure"
	"gopkg.in/yaml.v3"

	"github.com/spf13/viper"
)

type FrpsConfig struct {
	Host    string `json:"host" binding:"hostname"`        // Frps 服务器地址
	Port    int32  `json:"port" binding:"gte=0,lte=65535"` // Frps 服务器端口
	Token   string `json:"token"`                          // Frps 服务器 Token
	Allowed []struct {
		From    int32   `json:"from" binding:"gte=0,lte=65535"`         // Frps 服务器允许的端口范围
		To      int32   `json:"to" binding:"gte=0,lte=65535"`           // Frps 服务器允许的端口范围
		Exclude []int32 `json:"exclude" binding:"dive,gte=0,lte=65535"` // Frps 服务器排除的端口
	} `json:"allowed"`
}

type Config struct {
	Host string `mapstructure:"host" json:"host"` // 后端地址
	Path string `mapstructure:"path" json:"path"` // 数据存储路径
	Log  struct {
		Level string `mapstructure:"level" json:"level"` // 日志级别:DEBUG, INFO, WARNING, ERROR
		Save  bool   `mapstructure:"save" json:"save"`   // 是否保存日志到文件
	} `mapstructure:"log" json:"log"`
	AsyncQ struct {
		Log struct {
			Level string `mapstructure:"level" json:"level"`
		} `mapstructure:"log" json:"log"`
		Concurrency int `mapstructure:"concurrency" json:"concurrency"`
	} `mapstructure:"asynq" json:"asynq"`
	Gin struct {
		Mode   string `mapstructure:"mode" json:"mode"` // Gin 模式:debug, release, test
		Host   string `mapstructure:"host" json:"host"` // Gin 服务监听地址
		Port   uint   `mapstructure:"port" json:"port"` // Gin 服务监听端口
		Upload struct {
			Max int `mapstructure:"max" json:"max"` // 上传文件最大大小(单位:MB)
		} `mapstructure:"upload" json:"upload"`
		Proxies   []string `mapstructure:"proxies" json:"proxies"` // 信任的代理服务器
		RateLimit struct {
			Global    int      `mapstructure:"global" json:"global"`
			Whitelist []string `mapstructure:"whitelist" json:"whitelist"` // IP 白名单,不限制频率
		} `mapstructure:"ratelimit" json:"ratelimit"`
		CORS []string `mapstructure:"cors" json:"cors"`
		Log  struct {
			Whitelist []string `mapstructure:"whitelist" json:"whitelist"` // 日志白名单路径
		} `mapstructure:"log" json:"log"`
		JWT struct {
			Secret string `mapstructure:"secret" json:"secret"`
			Static bool   `mapstructure:"static" json:"static"`
		} `mapstructure:"jwt" json:"jwt"`
	} `mapstructure:"gin" json:"gin"`
	Gorm struct {
		MySQL struct {
			Host         string `mapstructure:"host" json:"host"`     // 数据库地址
			Port         uint   `mapstructure:"port" json:"port"`     // 数据库端口
			User         string `mapstructure:"user" json:"user"`     // 数据库用户名
			Pwd          string `mapstructure:"pwd" json:"pwd"`       // 数据库密码
			DB           string `mapstructure:"db" json:"db"`         // 数据库名称
			MaxOpenConns int    `mapstructure:"mxopen" json:"mxopen"` // 最大连接数
			MaxIdleConns int    `mapstructure:"mxidle" json:"mxidle"` // 最大空闲连接数
		} `mapstructure:"mysql" json:"mysql"` // MySQL 数据库配置
		Log struct {
			Level string `mapstructure:"level" json:"level"` // GORM 日志级别:INFO, WARNING, ERROR, SILENT
		} `mapstructure:"log" json:"log"`
	} `mapstructure:"gorm" json:"gorm"`
	Redis struct {
		Host string `mapstructure:"host" json:"host"` // Redis 地址
		Port uint   `mapstructure:"port" json:"port"` // Redis 端口
		Pwd  string `mapstructure:"pwd" json:"pwd"`   // Redis 密码
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
		Frp          struct {
			On         bool         `mapstructure:"on" json:"on"`
			FrpcImage  string       `mapstructure:"frpc" json:"frpc"`   // Frpc 镜像
			NginxImage string       `mapstructure:"nginx" json:"nginx"` // Nginx 镜像
			Frps       []FrpsConfig `mapstructure:"frps" json:"frps"`
		} `mapstructure:"frp" json:"frp"`
		GeneratorWorker int `mapstructure:"generator_worker" json:"generator_worker"`
	} `mapstructure:"k8s" json:"k8s"`
	NFS struct {
		Server  string `mapstructure:"server" json:"server"`
		Path    string `mapstructure:"path" json:"path"`
		Storage string `mapstructure:"storage" json:"storage"`
	} `mapstructure:"nfs" json:"nfs"`
	Cheat struct {
		IP struct {
			Whitelist []string `mapstructure:"whitelist" json:"whitelist"`
		} `mapstructure:"ip" json:"ip"`
	} `mapstructure:"cheat" json:"cheat"`
	GeoCityDB string `mapstructure:"geocity_db" json:"geocity_db"`
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
	Tidy()
}

// Tidy 格式化配置, 简单处理部分配置
func Tidy() {
	Env.Host = strings.TrimSuffix(Env.Host, "/")
	Env.Path = strings.TrimSuffix(Env.Path, "/")
	Env.NFS.Path = strings.TrimSuffix(Env.NFS.Path, "/")
}

// Save writes the current Env to config.yml.
// It preserves formatting only by overwriting the whole file.
func Save() error {
	if Env == nil {
		return errors.New("config env is nil")
	}
	settings := make(map[string]any)
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		TagName: "mapstructure",
		Result:  &settings,
	})
	if err != nil {
		return err
	}
	if err = decoder.Decode(Env); err != nil {
		return err
	}

	var root yaml.Node
	if err = yaml.Unmarshal(defaultConf, &root); err != nil {
		return err
	}
	if root.Kind == 0 {
		root.Kind = yaml.DocumentNode
	}
	if len(root.Content) == 0 {
		root.Content = []*yaml.Node{{Kind: yaml.MappingNode}}
	}
	if err = mergeYAMLNode(&root, settings); err != nil {
		return err
	}

	bytes, err := yaml.Marshal(&root)
	if err != nil {
		return err
	}
	return os.WriteFile("./config.yml", bytes, 0666)
}

func mergeYAMLNode(node *yaml.Node, data any) error {
	switch node.Kind {
	case yaml.DocumentNode:
		if len(node.Content) == 0 {
			node.Content = []*yaml.Node{{Kind: yaml.MappingNode}}
		}
		return mergeYAMLNode(node.Content[0], data)
	case yaml.MappingNode:
		mapped, ok := data.(map[string]any)
		if !ok {
			return nil
		}
		remaining := make(map[string]any, len(mapped))
		for key, value := range mapped {
			remaining[key] = value
		}
		for i := 0; i+1 < len(node.Content); i += 2 {
			keyNode := node.Content[i]
			valueNode := node.Content[i+1]
			if keyNode.Kind != yaml.ScalarNode {
				continue
			}
			key := keyNode.Value
			value, ok := remaining[key]
			if !ok {
				continue
			}
			if err := applyYAMLValue(valueNode, value); err != nil {
				return err
			}
			delete(remaining, key)
		}
		for key, value := range remaining {
			keyNode := &yaml.Node{Kind: yaml.ScalarNode, Value: key, Tag: "!!str"}
			valueNode, err := valueToNode(value)
			if err != nil {
				return err
			}
			node.Content = append(node.Content, keyNode, valueNode)
		}
		return nil
	case yaml.SequenceNode:
		items, ok := toSlice(data)
		if !ok {
			return nil
		}
		node.Content = node.Content[:0]
		for _, item := range items {
			itemNode, err := valueToNode(item)
			if err != nil {
				return err
			}
			node.Content = append(node.Content, itemNode)
		}
		return nil
	default:
		valueNode, err := valueToNode(data)
		if err != nil {
			return err
		}
		*node = *valueNode
		return nil
	}
}

func applyYAMLValue(node *yaml.Node, value any) error {
	switch node.Kind {
	case yaml.MappingNode, yaml.SequenceNode:
		return mergeYAMLNode(node, value)
	default:
		valueNode, err := valueToNode(value)
		if err != nil {
			return err
		}
		node.Kind = valueNode.Kind
		node.Tag = valueNode.Tag
		node.Value = valueNode.Value
		node.Content = valueNode.Content
		node.Style = valueNode.Style
		return nil
	}
}

func valueToNode(value any) (*yaml.Node, error) {
	bytes, err := yaml.Marshal(value)
	if err != nil {
		return nil, err
	}
	var node yaml.Node
	if err := yaml.Unmarshal(bytes, &node); err != nil {
		return nil, err
	}
	if node.Kind == yaml.DocumentNode && len(node.Content) > 0 {
		return node.Content[0], nil
	}
	return &node, nil
}

func toSlice(value any) ([]any, bool) {
	val := reflect.ValueOf(value)
	if !val.IsValid() {
		return nil, false
	}
	kind := val.Kind()
	if kind != reflect.Slice && kind != reflect.Array {
		return nil, false
	}
	items := make([]any, val.Len())
	for i := 0; i < val.Len(); i++ {
		items[i] = val.Index(i).Interface()
	}
	return items, true
}
