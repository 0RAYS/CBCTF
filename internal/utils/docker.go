package utils

import (
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"context"
	"time"

	"github.com/compose-spec/compose-go/loader"
	"github.com/compose-spec/compose-go/types"
	"gopkg.in/yaml.v3"
)

func LoadDockerComposeYaml(data string) (*types.Project, model.RetVal) {
	var raw map[string]any
	if err := yaml.Unmarshal([]byte(data), &raw); err != nil {
		log.Logger.Warningf("Failed to load docker-compose: %s", err)
		return nil, model.RetVal{Msg: i18n.Model.Docker.InvalidComposeYaml, Attr: map[string]any{"Error": err.Error()}}
	}
	if len(data) == 0 {
		return nil, model.RetVal{Msg: i18n.Model.Docker.InvalidComposeYaml, Attr: map[string]any{"Error": "Empty yaml"}}
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
		return nil, model.RetVal{Msg: i18n.Model.Docker.InvalidComposeYaml, Attr: map[string]any{"Error": err.Error()}}
	}
	return cfg, model.SuccessRetVal()
}
