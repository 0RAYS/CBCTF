package utils

import (
	"CBCTF/internel/i18n"
	"CBCTF/internel/log"
	"context"
	"github.com/compose-spec/compose-go/loader"
	"github.com/compose-spec/compose-go/types"
)

func LoadDockerComposeYaml(data string) (*types.Project, bool, string) {
	cfg, err := loader.LoadWithContext(context.Background(), types.ConfigDetails{
		ConfigFiles: []types.ConfigFile{
			{
				Content: []byte(data),
			},
		},
	})
	if err != nil {
		log.Logger.Warningf("Failed to load docker-compose.yml: %v", err)
		return nil, false, i18n.UnknownError
	}
	return cfg, true, i18n.Success
}
