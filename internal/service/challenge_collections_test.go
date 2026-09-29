package service

import (
	"testing"

	"CBCTF/internal/model"
)

func TestNativeVolumeCollectionKeepsComposeFlagsAndValidation(t *testing.T) {
	compose := `services:
  web:
    image: alpine:latest
    ports: ["80"]
    command: ["sleep", "60"]
    environment:
      EXAMPLE: kept
    x-volumes:
      - path: /flag
        content: static{volume-flag}
`
	template, flags, ret := buildChallengeTemplate(compose)
	if !ret.OK {
		t.Fatalf("compose load: %+v", ret)
	}
	if len(template.Pods) != 1 || len(template.Pods[0].Containers) != 1 {
		t.Fatalf("template: %+v", template)
	}
	container := template.Pods[0].Containers[0]
	if len(container.VolumeMounts) != 1 || container.VolumeMounts[0].Path != "/flag" || container.Environment["EXAMPLE"] != "kept" || len(container.Command) != 2 {
		t.Fatalf("native collections lost: %+v", container)
	}
	if len(flags) != 1 || flags[0].Value != "static{volume-flag}" || flags[0].Binding.Type != model.FileFlagBindingType {
		t.Fatalf("flag extension lost: %+v", flags)
	}
	_, _, ret = buildChallengeTemplate(compose + "      - path: /flag\n        content: duplicate\n")
	if ret.OK {
		t.Fatal("duplicate flag paths bypassed validation after collection change")
	}
}
