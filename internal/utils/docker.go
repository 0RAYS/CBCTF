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

func LoadDockerComposeYaml(data, project string, knownExtensions ...map[string]any) (*types.Project, error) {
	if strings.TrimSpace(data) == "" {
		return nil, errors.New("empty yaml")
	}
	var raw map[string]any
	if err := yaml.Unmarshal([]byte(data), &raw); err != nil {
		return nil, err
	}
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
		if len(knownExtensions) > 0 {
			options.KnownExtensions = knownExtensions[0]
		}
	})
}
