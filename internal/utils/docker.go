package utils

import (
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"context"
	"time"

	"github.com/compose-spec/compose-go/loader"
	"github.com/compose-spec/compose-go/types"
	"gopkg.in/yaml.v3"
)

func LoadDockerComposeYaml(data string) (*types.Project, bool, string) {
	var raw map[string]any
	if err := yaml.Unmarshal([]byte(data), &raw); err != nil || len(data) == 0 {
		log.Logger.Warningf("Failed to load docker-compose: %s", err)
		return nil, false, i18n.InvalidDockerComposeYaml
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	cfg, err := loader.LoadWithContext(ctx, types.ConfigDetails{
		ConfigFiles: []types.ConfigFile{
			{
				Config: raw,
			},
		},
	})
	if err != nil {
		log.Logger.Warningf("Failed to load docker-compose: %s", err)
		return nil, false, i18n.InvalidDockerComposeYaml
	}
	return cfg, true, i18n.Success
}
