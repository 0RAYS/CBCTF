package model

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"time"
)

// 靶场节点角色。
const (
	RangeNodeRole    = "node"
	RangeRouterRole  = "router"
	RangeGatewayRole = "gateway"
)

// DefaultRangeLeaseDuration 是作者未声明租约时长时的默认值。
const DefaultRangeLeaseDuration = 2 * time.Hour

// RangeConfig 靶场题目的调度与接入配置。
// PoolSize 是常驻就绪的实例数,租约结束即销毁实例并补建,不回收复用。
type RangeConfig struct {
	PoolSize      int           `json:"pool_size"`
	MaxLeased     int           `json:"max_leased"`
	LeaseDuration time.Duration `json:"lease_duration"`
	GatewayKey    string        `json:"gateway"`
	AccessNetwork string        `json:"access_network"`
	Routes        []RouteLink   `json:"routes"`
}

// RouteLink 声明两个网段之间三层互通,平台据此互相放行对方网段。
type RouteLink struct {
	From string `json:"from"`
	To   string `json:"to"`
}

// RangeAccess 是队伍登录跳板用的凭据。
// 公钥进跳板容器,私钥由平台保存并下发给队伍。
type RangeAccess struct {
	User      string
	PublicKey string
}

// RangeGateway 玩家网关的声明,来自出题人的 compose。
// 网关不属于实例拓扑,平台在租约时单独创建它,选手只拿到它的地址。
// 网关本身是平台固定的跳板容器,出题人只决定接入网段、静态地址与暴露端口。
type RangeGateway struct {
	Service string
	Name    string
	Ports   Exposes
	Network NetworkAttachment
	Subnet  NetworkDefinition
}

// GatewayOf 返回靶场题的玩家网关声明。
func GatewayOf(template ChallengeTemplate, accessNetwork string) (RangeGateway, bool) {
	for _, pod := range template.Pods {
		if pod.Role != RangeGatewayRole {
			continue
		}
		gateway := RangeGateway{
			Service: pod.Key,
			Name:    pod.Name,
			Ports:   pod.ServicePorts,
			Subnet:  NetworkDefinition{Name: accessNetwork},
		}
		for _, network := range pod.Networks {
			if network.Definition.Name == accessNetwork {
				gateway.Network = network.Attachment
				gateway.Subnet = network.Definition
			}
		}
		return gateway, true
	}
	return RangeGateway{}, false
}

// XRange 转换回 compose 声明,保证题目 YAML 往返不丢字段。
func (c RangeConfig) XRange() XRange {
	xRange := XRange{
		PoolSize:      c.PoolSize,
		MaxLeased:     c.MaxLeased,
		LeaseDuration: int(c.LeaseDuration / time.Second),
		Gateway:       c.GatewayKey,
		AccessNetwork: c.AccessNetwork,
	}
	for _, route := range c.Routes {
		xRange.Routes = append(xRange.Routes, XRangeRoute{From: route.From, To: route.To})
	}
	return xRange
}

// DeepCopy 返回靶场配置的副本。
func (c *RangeConfig) DeepCopy() *RangeConfig {
	if c == nil {
		return nil
	}
	copied := *c
	copied.Routes = append([]RouteLink(nil), c.Routes...)
	return &copied
}

// RoutedNetworks 返回与给定网段互通的网段名,互通关系是双向的。
func (c RangeConfig) RoutedNetworks(network string) []string {
	result := make([]string, 0)
	for _, route := range c.Routes {
		switch network {
		case route.From:
			result = append(result, route.To)
		case route.To:
			result = append(result, route.From)
		}
	}
	return result
}

// TemplateKey 计算渲染后拓扑的指纹:剔除每实例随机名后规范化哈希。
// 同一道靶场题的所有实例必须得到同一个指纹。
func TemplateKey(spec VictimSpec) string {
	type subnetKey struct {
		DefName    string   `json:"def_name"`
		CIDRBlock  string   `json:"cidr_block"`
		Gateway    string   `json:"gateway"`
		ExcludeIPs []string `json:"exclude_ips"`
	}
	normalized := struct {
		Pods    []PodSpec   `json:"pods"`
		Subnets []subnetKey `json:"subnets"`
	}{Pods: spec.Pods}
	for _, subnet := range spec.NetworkPlan.Subnets {
		if subnet == nil {
			continue
		}
		normalized.Subnets = append(normalized.Subnets, subnetKey{
			DefName:    subnet.DefName,
			CIDRBlock:  subnet.CIDRBlock,
			Gateway:    subnet.Gateway,
			ExcludeIPs: subnet.ExcludeIps,
		})
	}
	data, err := json.Marshal(normalized)
	if err != nil {
		return ""
	}
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}
