package config

import (
	_ "embed"
	"errors"
	"log"
	"os"
	"reflect"
	"strings"

	"github.com/google/uuid"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

type FrpsConfig struct {
	Host    string `json:"host" binding:"hostname"`
	Port    int32  `json:"port" binding:"gte=0,lte=65535"`
	Token   string `json:"token"`
	Allowed []struct {
		From    int32   `json:"from" binding:"gte=0,lte=65535"`
		To      int32   `json:"to" binding:"gte=0,lte=65535"`
		Exclude []int32 `json:"exclude" binding:"dive,gte=0,lte=65535"`
	} `json:"allowed"`
}

type Config struct {
	Host string `mapstructure:"host" json:"host"`
	Path string `mapstructure:"path" json:"path"`
	Log  struct {
		Level string `mapstructure:"level" json:"level"`
		Save  bool   `mapstructure:"save" json:"save"`
	} `mapstructure:"log" json:"log"`
	AsyncQ struct {
		Log struct {
			Level string `mapstructure:"level" json:"level"`
		} `mapstructure:"log" json:"log"`
		Concurrency int `mapstructure:"concurrency" json:"concurrency"`
		Queues      struct {
			Victim     int `mapstructure:"victim" json:"victim"`
			Generator  int `mapstructure:"generator" json:"generator"`
			Attachment int `mapstructure:"attachment" json:"attachment"`
			Email      int `mapstructure:"email" json:"email"`
			Webhook    int `mapstructure:"webhook" json:"webhook"`
			Image      int `mapstructure:"image" json:"image"`
		} `mapstructure:"queues" json:"queues"`
	} `mapstructure:"asynq" json:"asynq"`
	Gin struct {
		Mode   string `mapstructure:"mode" json:"mode"`
		Host   string `mapstructure:"host" json:"host"`
		Port   uint   `mapstructure:"port" json:"port"`
		Upload struct {
			Max int `mapstructure:"max" json:"max"`
		} `mapstructure:"upload" json:"upload"`
		Proxies   []string `mapstructure:"proxies" json:"proxies"`
		RateLimit struct {
			Global    int      `mapstructure:"global" json:"global"`
			Whitelist []string `mapstructure:"whitelist" json:"whitelist"`
		} `mapstructure:"ratelimit" json:"ratelimit"`
		CORS []string `mapstructure:"cors" json:"cors"`
		Log  struct {
			Whitelist []string `mapstructure:"whitelist" json:"whitelist"`
		} `mapstructure:"log" json:"log"`
		JWT struct {
			Secret string `mapstructure:"secret" json:"secret"`
		} `mapstructure:"jwt" json:"jwt"`
		Metrics struct {
			Whitelist []string `mapstructure:"whitelist" json:"whitelist"`
		} `mapstructure:"metrics" json:"metrics"`
	} `mapstructure:"gin" json:"gin"`
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
	K8S struct {
		Config       string `mapstructure:"config" json:"config"`
		Namespace    string `mapstructure:"namespace" json:"namespace"`
		TCPDumpImage string `mapstructure:"tcpdump" json:"tcpdump"`
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
	Registration struct {
		Enabled      bool `mapstructure:"enabled" json:"enabled"`
		DefaultGroup uint `mapstructure:"default_group" json:"default_group"`
	} `mapstructure:"registration" json:"registration"`
	GeoCityDB string `mapstructure:"geocity_db" json:"geocity_db"`
}

var Env *Config
var configFile string

//go:embed default.yaml
var defaultConf []byte

func Init(path string) {
	configFile = path
	viper.SetConfigFile(path)
	viper.SetEnvPrefix("CBCTF")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		if errors.Is(err, os.ErrNotExist) || errors.As(err, &viper.ConfigFileNotFoundError{}) {
			writeErr := os.WriteFile(path, defaultConf, 0600)
			if writeErr != nil {
				log.Fatalf("Failed to init config: %s", writeErr)
			}
			log.Fatalf("Config created at %s, please edit and restart", path)
		}
		log.Fatalf("Failed to read config: %s", err)
	}
	if err := viper.Unmarshal(&Env); err != nil {
		log.Fatalf("error unmarshalling config: %s", err)
	}
	tidy()
}

func tidy() {
	Env.Host = strings.TrimSuffix(Env.Host, "/")
	Env.Path = strings.TrimSuffix(Env.Path, "/")
	if Env.Gin.JWT.Secret == "" {
		Env.Gin.JWT.Secret = uuid.New().String()
	}
}

func Save() error {
	if Env == nil {
		return errors.New("config env is nil")
	}
	tidy()
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
	return os.WriteFile(configFile, bytes, 0600)
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
			remainingValue, exists := remaining[key]
			if !exists {
				continue
			}
			if err := func(node *yaml.Node, value any) error {
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
			}(valueNode, remainingValue); err != nil {
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
		items, ok := func(value any) ([]any, bool) {
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
		}(data)
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

func valueToNode(value any) (*yaml.Node, error) {
	bytes, err := yaml.Marshal(value)
	if err != nil {
		return nil, err
	}
	var node yaml.Node
	if err = yaml.Unmarshal(bytes, &node); err != nil {
		return nil, err
	}
	if node.Kind == yaml.DocumentNode && len(node.Content) > 0 {
		return node.Content[0], nil
	}
	return &node, nil
}
