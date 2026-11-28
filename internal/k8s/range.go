package k8s

import (
	"context"
	"strings"
	"time"

	corev1 "k8s.io/api/core/v1"

	"CBCTF/internal/config"
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
)

// rangeGatewayLabel 区分玩家跳板与拓扑节点,Service 只选中跳板。
const rangeGatewayLabel = "range_gateway"

// StartRangeGateway 创建玩家跳板。
// 跳板是租约级资源:SSH 配置、Pod 与 Service 都挂在实例的资源树上,实例销毁时一并回收。
// 跳板是平台固定容器,只提供 SSH 登录,题目相关的环境变量、flag 与文件都不进这个容器。
func StartRangeGateway(ctx context.Context, victim model.Victim, gateway model.RangeGateway, access model.RangeAccess) (model.Victim, model.RetVal) {
	subnet := rangeSubnet(victim, gateway.Subnet.Name)
	if subnet == nil {
		return victim, model.FailRetVal(i18n.K8S.GetError, "Subnet", gateway.Subnet.Name, errMissingSubnet)
	}
	sshPort, ret := gatewaySSHPort(gateway)
	if !ret.OK {
		return victim, ret
	}
	ctx, ret = withVictimRootOwner(ctx, victim, "workloads")
	if !ret.OK {
		return victim, ret
	}
	labels := VictimLabels(victim, map[string]string{RoleLabel: VictimPodTag, rangeGatewayLabel: "true"})
	podName := victimResourceName("rgt", victim.ID, "gateway")
	sshConfig, ret := CreateConfigMap(ctx, CreateConfigMapOptions{
		Name:   victimResourceName("rpay", victim.ID, "gateway"),
		Labels: labels,
		Data:   jumpSSHConfig(access.User, access.PublicKey),
	})
	if !ret.OK {
		return victim, ret
	}
	if sshConfig == nil {
		return victim, model.FailRetVal(i18n.K8S.CreateError, "ConfigMap", sshConfig.Name, errMissingConfigMap)
	}

	containers := []corev1.Container{jumpContainer(sshPort)}
	volumes := []corev1.Volume{{
		Name:      jumpSSHVolume,
		ConfigMap: &corev1.ConfigMapVolumeSource{Name: sshConfig.Name},
	}}
	if config.Env.K8S.CaptureEnabled {
		// The capture sidecar writes into the traffic directory, which has to
		// exist and belong to the capture user before the pod is created.
		if err := prepareCaptureTrafficDir(victim); err != nil {
			log.Logger.Warningf("Failed to prepare traffic directory for victim %d: %s", victim.ID, err)
		}
		containers = append(containers, captureContainer(victim, "jump.pcap"))
		volumes = append(volumes, nfsVolume())
	}
	frpSidecar, ret := jumpFrpSidecars(ctx, victim, labels, gateway, sshPort)
	if !ret.OK {
		return victim, ret
	}
	containers = append(containers, frpSidecar.Containers...)
	volumes = append(volumes, frpSidecar.Volumes...)

	pod, ret := CreatePod(ctx, CreatePodOptions{
		SubmitOnly:         true,
		Name:               podName,
		Labels:             labels,
		Containers:         containers,
		Volumes:            volumes,
		ServiceAccountName: runtimeServiceAccount(),
		Networks: []Network{{
			Interface:    gateway.Network.Name,
			IPv4:         gateway.Network.IP,
			MAC:          gateway.Network.MAC,
			CIDR:         gateway.Subnet.CIDR,
			Subnet:       subnet.Name,
			NetAttachDef: subnet.NetAttachDef.Name,
			Routes:       rangeRoutes(victim, gateway.Subnet.Name),
		}},
	})
	if !ret.OK {
		return victim, ret
	}
	if pod == nil {
		return victim, model.FailRetVal(i18n.K8S.CreateError, "Pod", podName, errMissingPodObject)
	}
	victim.Resources.UIDs[podName] = string(pod.UID)
	victim.Pods = append(victim.Pods, model.Pod{VictimID: victim.ID, Name: podName})

	if ret = createRangeGatewayPolicy(ctx, victim, gateway, frpSidecar.Endpoints); !ret.OK {
		return victim, ret
	}

	service, ret := CreateService(ctx, CreateServiceOptions{
		Name:     victimResourceName("rsvc", victim.ID, "gateway"),
		Labels:   labels,
		Ports:    gateway.Ports,
		Selector: labels,
	})
	if !ret.OK {
		return victim, ret
	}
	if service == nil {
		return victim, model.FailRetVal(i18n.K8S.CreateError, "Service", service.Name, errMissingServiceObject)
	}
	hostIP, ret := waitPodHostIP(ctx, podName, 30*time.Second)
	if !ret.OK {
		return victim, ret
	}
	for _, port := range service.Spec.Ports {
		victim.Endpoints = append(victim.Endpoints, model.Endpoint{
			Name:     findExposeDisplayName(gateway.Ports, port.Port, string(port.Protocol)),
			IP:       hostIP,
			Port:     port.NodePort,
			Protocol: string(port.Protocol),
		})
	}
	// NodePort 与 frps 两个入口并存,选手拿到哪个都能进跳板。
	victim.ExposedEndpoints = append(model.Endpoints(nil), victim.Endpoints...)
	victim.ExposedEndpoints = append(victim.ExposedEndpoints, frpSidecar.Endpoints...)
	log.Logger.Infof(
		"Range gateway created: victim_id=%d team_id=%d pod=%s endpoints=%d",
		victim.ID, victim.TeamID.V, podName, len(victim.ExposedEndpoints),
	)
	return victim, model.SuccessRetVal()
}

// gatewaySSHPort 返回跳板对外暴露的 SSH 端口。
func gatewaySSHPort(gateway model.RangeGateway) (int32, model.RetVal) {
	for _, port := range gateway.Ports {
		if strings.EqualFold(port.Protocol, "tcp") {
			return port.Port, model.SuccessRetVal()
		}
	}
	return 0, model.FailRetVal(i18n.K8S.CreateError, "RangeGateway", "ssh port", nil)
}

// jumpFrpSidecars 给跳板注入 frpc 与 nginx 侧车,SSH 端口同时经 frps 对外暴露。
func jumpFrpSidecars(ctx context.Context, victim model.Victim, labels map[string]string, gateway model.RangeGateway, sshPort int32) (frpSidecars, model.RetVal) {
	if !config.Env.K8S.Frp.On || len(config.Env.K8S.Frp.Frps) == 0 {
		return frpSidecars{}, model.SuccessRetVal()
	}
	sidecars, ret := buildFrpSidecars(victimResourceName("rfrp", victim.ID, "gateway"), []model.Endpoint{{
		Name:     findExposeDisplayName(gateway.Ports, sshPort, "tcp"),
		IP:       "127.0.0.1",
		Port:     sshPort,
		Protocol: "tcp",
	}})
	if !ret.OK {
		return frpSidecars{}, ret
	}
	frpConfig, ret := CreateConfigMap(ctx, CreateConfigMapOptions{
		Name:   sidecars.ConfigName,
		Labels: labels,
		Data:   sidecars.ConfigData,
	})
	if !ret.OK {
		return frpSidecars{}, ret
	}
	if frpConfig == nil {
		return frpSidecars{}, model.FailRetVal(i18n.K8S.CreateError, "ConfigMap", sidecars.ConfigName, errMissingConfigMap)
	}
	return sidecars, model.SuccessRetVal()
}

// rangeSubnet 返回实例拓扑中指定网段的子网。
func rangeSubnet(victim model.Victim, name string) *model.Subnet {
	for _, subnet := range victim.Spec.NetworkPlan.Subnets {
		if subnet != nil && subnet.DefName == name {
			return subnet
		}
	}
	return nil
}

// createRangeGatewayPolicy 给玩家跳板建出站策略。
// 拓扑节点的策略按 service 选择器绑定,跳板在拓扑之外,
// 没有策略的 Pod 出站不受限制,跳板因此要自己带上放行范围。
func createRangeGatewayPolicy(ctx context.Context, victim model.Victim, gateway model.RangeGateway, frpEndpoints model.Endpoints) model.RetVal {
	labels := VictimLabels(victim, map[string]string{rangeGatewayLabel: "true"})
	selector := map[string]string{
		"victim_id":       labels["victim_id"],
		rangeGatewayLabel: "true",
	}
	routedCIDRs := make([]string, 0)
	for _, link := range victim.Spec.Routes {
		peer := ""
		switch gateway.Subnet.Name {
		case link.From:
			peer = link.To
		case link.To:
			peer = link.From
		default:
			continue
		}
		if subnet := rangeSubnet(victim, peer); subnet != nil {
			routedCIDRs = append(routedCIDRs, subnet.CIDRBlock)
		}
	}
	peerIPs := make(map[string][]string)
	for _, podSpec := range victim.Spec.Pods {
		service := podServiceName(podSpec)
		if service == "" {
			continue
		}
		ips := make([]string, 0, len(podSpec.Networks))
		for _, network := range podSpec.Networks {
			if network.Attachment.IP != "" {
				ips = append(ips, network.Attachment.IP)
			}
		}
		if len(ips) > 0 {
			peerIPs[service] = ips
		}
	}
	// NetworkPolicy 按 Pod 生效,注入跳板的 frpc 没法靠标签躲开出站策略,
	// frps 服务器的地址要单独放行。
	egressPeers := make([]string, 0)
	for _, endpoint := range frpEndpoints {
		egressPeers = append(egressPeers, frpsPeerCIDRs(endpoint.IP)...)
	}
	_, ret := CreateNetworkPolicy(ctx, CreateNetworkPolicyOptions{
		Name:        victimResourceName("npg", victim.ID, "gateway"),
		Labels:      labels,
		Selector:    selector,
		Service:     gateway.Service,
		Policies:    victim.Spec.NetworkPolicies,
		SubnetCIDRs: []string{gateway.Subnet.CIDR},
		PeerIPs:     peerIPs,
		VPCSubnets:  vpcSubnetCIDRs(victim),
		RoutedCIDRs: routedCIDRs,
		EgressPeers: egressPeers,
	})
	return ret
}

// rangeRoutes 返回节点到对端网段的静态路由。
// 节点只连着本网段,到互通网段的流量交给本网段的网关转发;
// 集群默认路由不属于靶场,不能用它跨段。
func rangeRoutes(victim model.Victim, network string) []model.Route {
	cidrByName := make(map[string]string)
	gateway := ""
	for _, subnet := range victim.Spec.NetworkPlan.Subnets {
		if subnet == nil {
			continue
		}
		cidrByName[subnet.DefName] = subnet.CIDRBlock
		if subnet.DefName == network {
			gateway = subnet.Gateway
		}
	}
	if gateway == "" {
		return nil
	}
	routes := make([]model.Route, 0)
	for _, link := range victim.Spec.Routes {
		peer := ""
		switch network {
		case link.From:
			peer = link.To
		case link.To:
			peer = link.From
		default:
			continue
		}
		if cidr := cidrByName[peer]; cidr != "" {
			routes = append(routes, model.Route{Dst: cidr, Gw: gateway})
		}
	}
	return routes
}

// waitPodHostIP 等待 Pod 被调度到节点上,NodePort 的入口地址依赖它。
func waitPodHostIP(ctx context.Context, name string, timeout time.Duration) (string, model.RetVal) {
	deadline := time.Now().Add(timeout)
	for {
		pod, ret := GetPod(ctx, name)
		if !ret.OK {
			return "", ret
		}
		if pod.Status.HostIP != "" {
			return pod.Status.HostIP, model.SuccessRetVal()
		}
		if time.Now().After(deadline) {
			return "", model.FailRetVal(i18n.K8S.PodRunError, "Pod", name, errPodHostIPTimeout)
		}
		select {
		case <-ctx.Done():
			return "", model.FailRetVal(i18n.Common.UnknownError, "Pod", name, ctx.Err())
		case <-time.After(500 * time.Millisecond):
		}
	}
}
