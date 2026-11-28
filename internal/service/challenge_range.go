package service

import (
	"fmt"
	"strings"

	"github.com/compose-spec/compose-go/v2/types"

	"CBCTF/internal/config"
	"CBCTF/internal/model"
)

// rangeChallenge 判断 compose 是否声明了靶场拓扑。
func rangeChallenge(project *types.Project) bool {
	_, ok := project.Extensions[model.XRangeExtension].(model.XRange)
	return ok
}

// definedNetworkNames 返回 compose 中定义的自定义网段名。
// compose 会给网段名加上项目名前缀,作者在 x-range 里写的是不带前缀的名字。
func definedNetworkNames(project *types.Project) map[string]struct{} {
	names := make(map[string]struct{})
	for _, network := range project.Networks {
		name := strings.TrimPrefix(network.Name, project.Name+"_")
		if name == "" || name == "default" {
			continue
		}
		names[name] = struct{}{}
	}
	return names
}

// parseRangeTopology 解析 x-range 声明并校验拓扑约束。
// 玩家网关是选手的唯一入口,由平台在租约时创建,不进入实例拓扑;
// 网段互通是双向的,平台据此互相放行对方网段。
func parseRangeTopology(project *types.Project) (*model.RangeConfig, model.RetVal) {
	xRange, ok := project.Extensions[model.XRangeExtension].(model.XRange)
	if !ok {
		return nil, model.SuccessRetVal()
	}
	rangeConfig := xRange.RangeConfig()
	if rangeConfig.PoolSize < 0 || rangeConfig.MaxLeased < 0 {
		return nil, invalidComposeYamlRetVal("x-range pool_size and max_leased must not be negative")
	}
	if rangeConfig.LeaseDuration <= 0 {
		rangeConfig.LeaseDuration = model.DefaultRangeLeaseDuration
	}
	networks := definedNetworkNames(project)
	gateway, ok := project.Services[rangeConfig.GatewayKey]
	if !ok {
		return nil, invalidComposeYamlRetVal(fmt.Sprintf("x-range gateway %s is not a defined service", rangeConfig.GatewayKey))
	}
	if kubeVirt, _ := gateway.Extensions[model.XKubeVirtExtension].(bool); kubeVirt {
		return nil, invalidComposeYamlRetVal(fmt.Sprintf("x-range gateway %s must not run as a virtual machine", rangeConfig.GatewayKey))
	}
	if rangeConfig.AccessNetwork == "" {
		if len(gateway.Networks) != 1 {
			return nil, invalidComposeYamlRetVal("x-range access_network is required when the gateway has more than one network")
		}
		for name := range gateway.Networks {
			rangeConfig.AccessNetwork = name
		}
	}
	if _, ok := gateway.Networks[rangeConfig.AccessNetwork]; !ok {
		return nil, invalidComposeYamlRetVal(fmt.Sprintf("x-range gateway %s is not attached to access network %s", rangeConfig.GatewayKey, rangeConfig.AccessNetwork))
	}
	if _, ok := networks[rangeConfig.AccessNetwork]; !ok {
		return nil, invalidComposeYamlRetVal(fmt.Sprintf("x-range access network %s is not defined", rangeConfig.AccessNetwork))
	}
	seen := make(map[string]struct{})
	for _, route := range rangeConfig.Routes {
		if route.From == route.To {
			return nil, invalidComposeYamlRetVal(fmt.Sprintf("x-range route %s to %s must connect two different networks", route.From, route.To))
		}
		for _, name := range []string{route.From, route.To} {
			if _, ok := networks[name]; !ok {
				return nil, invalidComposeYamlRetVal(fmt.Sprintf("x-range route network %s is not defined", name))
			}
		}
		key := route.From + "|" + route.To
		if _, ok := seen[key]; ok {
			return nil, invalidComposeYamlRetVal(fmt.Sprintf("x-range route %s to %s is duplicated", route.From, route.To))
		}
		if _, ok := seen[route.To+"|"+route.From]; ok {
			return nil, invalidComposeYamlRetVal(fmt.Sprintf("x-range route %s to %s is duplicated", route.From, route.To))
		}
		seen[key] = struct{}{}
	}
	// 选手从接入网段进入拓扑,接入网段不与任何网段互通时选手够不到别的网段上的节点。
	if !rangeAccessRouted(rangeConfig) {
		for _, service := range project.Services {
			if service.Name == rangeConfig.GatewayKey {
				continue
			}
			for name := range service.Networks {
				if name != rangeConfig.AccessNetwork {
					return nil, invalidComposeYamlRetVal(fmt.Sprintf("x-range access network %s must be connected by a route", rangeConfig.AccessNetwork))
				}
			}
		}
	}
	for _, service := range project.Services {
		role, _ := service.Extensions[model.XRoleExtension].(string)
		switch role {
		case "", model.RangeNodeRole, model.RangeRouterRole:
		case model.RangeGatewayRole:
			if service.Name != rangeConfig.GatewayKey {
				return nil, invalidComposeYamlRetVal(fmt.Sprintf("service %s has the gateway role but x-range gateway is %s", service.Name, rangeConfig.GatewayKey))
			}
		default:
			return nil, invalidComposeYamlRetVal(fmt.Sprintf("service %s has an unknown role %s", service.Name, role))
		}
	}
	return &rangeConfig, model.SuccessRetVal()
}

// rangeAccessRouted 判断玩家接入网段是否与其它网段互通。
func rangeAccessRouted(rangeConfig model.RangeConfig) bool {
	for _, route := range rangeConfig.Routes {
		if route.From == rangeConfig.AccessNetwork || route.To == rangeConfig.AccessNetwork {
			return true
		}
	}
	return false
}

// fillPlatformContainerImages 给平台固定容器补上镜像。
// 玩家网关是平台固定的跳板容器,作者只写需要配置的字段,镜像由平台提供。
func fillPlatformContainerImages(raw map[string]any) error {
	services, _ := raw["services"].(map[string]any)
	if services == nil {
		return nil
	}
	gateway := ""
	if xRange, ok := raw[model.XRangeExtension].(map[string]any); ok {
		gateway, _ = xRange["gateway"].(string)
	}
	for name, value := range services {
		service, _ := value.(map[string]any)
		if service == nil {
			continue
		}
		role, _ := service[model.XRoleExtension].(string)
		if name != gateway && role != model.RangeGatewayRole {
			continue
		}
		if image, _ := service["image"].(string); strings.TrimSpace(image) != "" {
			continue
		}
		if config.Env.K8S.Range.JumperImage == "" {
			return fmt.Errorf("the player gateway image is not configured")
		}
		service["image"] = config.Env.K8S.Range.JumperImage
	}
	return nil
}

// rangeRole 返回 service 在靶场拓扑中的角色。
func rangeRole(rangeConfig *model.RangeConfig, service types.ServiceConfig) string {
	if rangeConfig == nil {
		return ""
	}
	if service.Name == rangeConfig.GatewayKey {
		return model.RangeGatewayRole
	}
	if role, _ := service.Extensions[model.XRoleExtension].(string); role != "" {
		return role
	}
	return model.RangeNodeRole
}
