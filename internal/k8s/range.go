package k8s

import (
	"context"
	"encoding/json"
	"fmt"
	"maps"
	"slices"
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"

	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
)

// rangeGatewayPayloadPath 是玩家网关读取租约载荷的挂载点。
const rangeGatewayPayloadPath = "/run/range/payload"

// rangeGatewayLabel 区分玩家网关与拓扑节点,Service 只选中网关。
const rangeGatewayLabel = "range_gateway"

// StartRangeGateway 创建玩家网关。
// 网关是租约级资源:载荷、Pod 与 Service 都挂在实例的资源树上,实例销毁时一并回收。
// 选手只拿到网关的地址,拓扑节点不对外暴露。
func StartRangeGateway(ctx context.Context, victim model.Victim, gateway model.RangeGateway, payload model.GatewayPayload) (model.Victim, model.RetVal) {
	subnet := rangeSubnet(victim, gateway.Subnet.Name)
	if subnet == nil {
		return victim, model.FailRetVal(i18n.K8S.GetError, "Subnet", gateway.Subnet.Name, errMissingSubnet)
	}
	labels := VictimLabels(victim, map[string]string{RoleLabel: VictimPodTag, rangeGatewayLabel: "true"})
	podName := victimResourceName("rgt", victim.ID, "gateway")

	data := make(map[string]string, len(payload.Files)+1)
	maps.Copy(data, payload.Files)
	meta, err := json.Marshal(map[string]any{
		"credentials": payload.Credentials,
		"env":         payload.Env,
		"files":       slices.Sorted(maps.Keys(payload.Files)),
	})
	if err != nil {
		return victim, model.FailRetVal(i18n.Common.UnknownError, "GatewayPayload", podName, err)
	}
	data["payload.json"] = string(meta)
	cm, ret := CreateConfigMap(ctx, CreateConfigMapOptions{
		Name:   victimResourceName("rpay", victim.ID, "gateway"),
		Labels: labels,
		Data:   data,
	})
	if !ret.OK {
		return victim, ret
	}
	if cm == nil {
		return victim, model.FailRetVal(i18n.K8S.CreateError, "ConfigMap", cm.Name, errMissingConfigMap)
	}

	envs := make([]corev1.EnvVar, 0, len(payload.Env))
	for key, value := range payload.Env {
		envs = append(envs, corev1.EnvVar{Name: key, Value: value})
	}
	slices.SortFunc(envs, func(a, b corev1.EnvVar) int {
		if a.Name < b.Name {
			return -1
		}
		return 1
	})
	limit := corev1.ResourceList{}
	if gateway.Resources.CPUMillis > 0 {
		limit[corev1.ResourceCPU] = resource.MustParse(fmt.Sprintf("%dm", gateway.Resources.CPUMillis))
	}
	if gateway.Resources.MemoryBytes > 0 {
		limit[corev1.ResourceMemory] = resource.MustParse(fmt.Sprintf("%d", gateway.Resources.MemoryBytes))
	}
	container := corev1.Container{
		Name:            gateway.Name,
		Image:           gateway.Image,
		ImagePullPolicy: corev1.PullIfNotPresent,
		Env:             envs,
		Command:         gateway.Command,
		WorkingDir:      gateway.WorkingDir,
		Resources: corev1.ResourceRequirements{
			Limits:   limit,
			Requests: workloadRequests(limit),
		},
		VolumeMounts: []corev1.VolumeMount{{
			Name:      "range-payload",
			MountPath: rangeGatewayPayloadPath,
			ReadOnly:  true,
		}},
	}
	volumes := []corev1.Volume{{
		Name:      "range-payload",
		ConfigMap: &corev1.ConfigMapVolumeSource{Name: cm.Name},
	}}
	pod, ret := CreatePod(ctx, CreatePodOptions{
		SubmitOnly:         true,
		Name:               podName,
		Labels:             labels,
		Containers:         []corev1.Container{container},
		Volumes:            volumes,
		ServiceAccountName: runtimeServiceAccount(),
		Networks: []Network{{
			Interface:    gateway.Network.Name,
			IPv4:         gateway.Network.IP,
			MAC:          gateway.Network.MAC,
			CIDR:         gateway.Subnet.CIDR,
			Gateway:      gateway.Subnet.Gateway,
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
	victim.ExposedEndpoints = append(model.Endpoints(nil), victim.Endpoints...)
	log.Logger.Infof(
		"Range gateway created: victim_id=%d team_id=%d pod=%s endpoints=%d",
		victim.ID, victim.TeamID.V, podName, len(victim.ExposedEndpoints),
	)
	return victim, model.SuccessRetVal()
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
