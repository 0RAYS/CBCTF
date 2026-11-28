package utils

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/compose-spec/compose-go/v2/loader"
	"github.com/compose-spec/compose-go/v2/types"
	"go.yaml.in/yaml/v4"
)

// ParseComposeYaml 解析 compose 原文,调用方可以在加载前补齐平台固定容器的镜像。
func ParseComposeYaml(data string) (map[string]any, error) {
	if strings.TrimSpace(data) == "" {
		return nil, errors.New("empty yaml")
	}
	var raw map[string]any
	if err := yaml.Unmarshal([]byte(data), &raw); err != nil {
		return nil, err
	}
	return raw, nil
}

func LoadComposeProject(raw map[string]any, project string, knownExtensions map[string]any) (*types.Project, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	return loader.LoadWithContext(ctx, types.ConfigDetails{
		ConfigFiles: []types.ConfigFile{
			{
				Filename: "-",
				Config:   raw,
			},
		},
	}, func(options *loader.Options) {
		options.SetProjectName(project, true)
		if knownExtensions != nil {
			options.KnownExtensions = knownExtensions
		}
	})
}
