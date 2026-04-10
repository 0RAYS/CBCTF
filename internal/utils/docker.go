package utils

import (
	"context"
	"errors"
	"time"

	"github.com/compose-spec/compose-go/loader"
	"github.com/compose-spec/compose-go/types"
	"gopkg.in/yaml.v3"
)

func LoadDockerComposeYaml(data string) (*types.Project, error) {
	var raw map[string]any
	if err := yaml.Unmarshal([]byte(data), &raw); err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, errors.New("empty yaml")
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	return loader.LoadWithContext(ctx, types.ConfigDetails{
		ConfigFiles: []types.ConfigFile{
			{
				Config: raw,
			},
		},
	})
}
