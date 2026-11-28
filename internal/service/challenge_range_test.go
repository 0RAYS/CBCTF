package service

import (
	"strings"
	"testing"
	"time"

	"CBCTF/internal/model"
)

const rangeCompose = `x-range:
  pool_size: 3
  max_leased: 10
  lease_duration: 1800
  gateway: jump
  access_network: mgmt
  routes:
    - {from: mgmt, to: internal}
    - {from: internal, to: dmz}
services:
  jump:
    image: alpine:latest
    x-role: gateway
    ports:
      - {target: 22, published: ssh}
    networks:
      mgmt: {ipv4_address: 192.168.99.254}
  router:
    image: alpine:latest
    x-role: router
    networks:
      internal: {ipv4_address: 10.10.0.2}
      dmz: {ipv4_address: 10.20.0.2}
  web:
    image: nginx:latest
    networks:
      internal: {ipv4_address: 10.10.0.3}
networks:
  mgmt:
    ipam:
      config:
        - subnet: 192.168.99.0/24
          gateway: 192.168.99.1
  internal:
    ipam:
      config:
        - subnet: 10.10.0.0/24
          gateway: 10.10.0.1
  dmz:
    ipam:
      config:
        - subnet: 10.20.0.0/24
          gateway: 10.20.0.1
`

func TestRangeTopologyParsesExtension(t *testing.T) {
	template, _, rangeConfig, ret := buildChallengeTemplate(rangeCompose)
	if !ret.OK {
		t.Fatalf("compose load: %+v", ret)
	}
	if rangeConfig == nil {
		t.Fatal("range config missing")
	}
	if rangeConfig.PoolSize != 3 || rangeConfig.MaxLeased != 10 || rangeConfig.GatewayKey != "jump" || rangeConfig.AccessNetwork != "mgmt" {
		t.Fatalf("range config: %+v", rangeConfig)
	}
	if rangeConfig.LeaseDuration != 30*time.Minute {
		t.Fatalf("lease duration: %v", rangeConfig.LeaseDuration)
	}
	if len(rangeConfig.Routes) != 2 || rangeConfig.Routes[0].From != "mgmt" || rangeConfig.Routes[0].To != "internal" {
		t.Fatalf("routes: %+v", rangeConfig.Routes)
	}
	roles := make(map[string]string)
	for _, pod := range template.Pods {
		roles[pod.Name] = pod.Role
	}
	if roles["jump"] != model.RangeGatewayRole || roles["router"] != model.RangeRouterRole || roles["web"] != model.RangeNodeRole {
		t.Fatalf("roles: %+v", roles)
	}
}

func TestRangeTopologyRejectsInvalidDeclarations(t *testing.T) {
	cases := map[string]string{
		"unknown gateway": strings.Replace(rangeCompose, "gateway: jump", "gateway: missing", 1),
		"unknown route":   strings.Replace(rangeCompose, "{from: internal, to: dmz}", "{from: internal, to: missing}", 1),
		"self route":      strings.Replace(rangeCompose, "{from: internal, to: dmz}", "{from: internal, to: internal}", 1),
		"role without range": strings.Replace(
			strings.Replace(rangeCompose, "x-range:\n  pool_size: 3\n  max_leased: 10\n  lease_duration: 1800\n  gateway: jump\n  access_network: mgmt\n  routes:\n    - {from: mgmt, to: internal}\n    - {from: internal, to: dmz}\n", "", 1),
			"x-role: router", "", 1),
		"access network without route": strings.Replace(
			strings.Replace(rangeCompose, "    - {from: mgmt, to: internal}\n", "", 1),
			"access_network: mgmt", "access_network: nowhere", 1),
	}
	for name, compose := range cases {
		t.Run(name, func(t *testing.T) {
			if _, _, _, ret := buildChallengeTemplate(compose); ret.OK {
				t.Fatalf("compose accepted: %s", name)
			}
		})
	}
}

func TestRangeTopologySurvivesYamlRoundTrip(t *testing.T) {
	template, flags, rangeConfig, ret := buildChallengeTemplate(rangeCompose)
	if !ret.OK {
		t.Fatalf("compose load: %+v", ret)
	}
	yaml := Template2Yaml(template, flags, rangeConfig)
	if !strings.Contains(yaml, "x-range:") || !strings.Contains(yaml, "x-role: gateway") {
		t.Fatalf("range extensions lost: %s", yaml)
	}
	_, _, parsed, ret := buildChallengeTemplate(yaml)
	if !ret.OK {
		t.Fatalf("round trip load: %+v", ret)
	}
	if parsed == nil || parsed.GatewayKey != "jump" || len(parsed.Routes) != 2 {
		t.Fatalf("round trip config: %+v", parsed)
	}
}
