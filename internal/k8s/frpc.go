package k8s

import (
	"bytes"
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"net"
	"slices"
	"strings"

	"github.com/BurntSushi/toml"
	corev1 "k8s.io/api/core/v1"

	"CBCTF/internal/config"
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"CBCTF/internal/redis"
	"CBCTF/internal/utils"
)

type frpcConfig struct {
	ServerAddr string      `toml:"serverAddr"`
	Auth       frpcAuth    `toml:"auth"`
	Proxies    []frpcProxy `toml:"proxies"`
	ServerPort int32       `toml:"serverPort"`
}

type frpcAuth struct {
	Token string `toml:"token"`
}

type frpcProxy struct {
	Transport  *frpcTransport `toml:"transport,omitempty"`
	Name       string         `toml:"name"`
	Type       string         `toml:"type"`
	LocalIP    string         `toml:"localIP"`
	LocalPort  int32          `toml:"localPort"`
	RemotePort int32          `toml:"remotePort"`
}

type frpcTransport struct {
	ProxyProtocolVersion string `toml:"proxyProtocolVersion"`
}

type nginxConfig struct {
	Streams []nginxStream
}

type nginxStream struct {
	Name       string
	TargetIP   string
	TargetPort int32
	ListenPort int32
}

func (c frpcConfig) String() (string, error) {
	var buf bytes.Buffer
	if err := toml.NewEncoder(&buf).Encode(c); err != nil {
		return "", err
	}
	return strings.TrimSpace(buf.String()), nil
}

func (c nginxConfig) String() string {
	var buf bytes.Buffer
	buf.WriteString("worker_processes auto;\n")
	buf.WriteString("events {\n")
	buf.WriteString("    worker_connections 1024;\n")
	buf.WriteString("}\n")
	buf.WriteString("stream {\n")
	for _, stream := range c.Streams {
		server := net.JoinHostPort(stream.TargetIP, fmt.Sprintf("%d", stream.TargetPort))
		_, _ = fmt.Fprintf(&buf, "    upstream %s {\n", stream.Name)
		_, _ = fmt.Fprintf(&buf, "        server %s;\n", server)
		buf.WriteString("    }\n")
		buf.WriteString("    server {\n")
		_, _ = fmt.Fprintf(&buf, "        listen %d proxy_protocol;\n", stream.ListenPort)
		_, _ = fmt.Fprintf(&buf, "        proxy_pass %s;\n", stream.Name)
		buf.WriteString("    }\n")
	}
	buf.WriteString("}")
	return buf.String()
}

func addFrpcProxy(config *frpcConfig, protocol string, localIP string, localPort int32, remotePort int32, proxyProtocol bool) {
	proxy := frpcProxy{
		Name:       utils.RandHexStr(10),
		Type:       protocol,
		LocalIP:    localIP,
		LocalPort:  localPort,
		RemotePort: remotePort,
	}
	if proxyProtocol {
		proxy.Transport = &frpcTransport{ProxyProtocolVersion: "v2"}
	}
	config.Proxies = append(config.Proxies, proxy)
}

// errNoFrpsServer marks a victim that cannot be exposed because no FRPS
// server is configured.
var errNoFrpsServer = errors.New("no FRPS server is configured")

// frpSidecars 描述注入 Pod 的 frpc 与 nginx 侧车。
type frpSidecars struct {
	Containers []corev1.Container
	Volumes    []corev1.Volume
	ConfigName string
	ConfigData map[string]string
	Endpoints  model.Endpoints
}

// buildFrpSidecars 生成 frpc 与 nginx 侧车。
// targets 是经 frps 暴露的目标,侧车注入到别的 Pod 里时只能填 Pod 内的地址。
func buildFrpSidecars(configName string, targets []model.Endpoint) (frpSidecars, model.RetVal) {
	frps, ret := pickFrps()
	if !ret.OK {
		return frpSidecars{}, ret
	}
	portRange := frpsPortRange(frps)
	frpc := frpcConfig{
		ServerAddr: frps.Host,
		ServerPort: frps.Port,
		Auth:       frpcAuth{Token: frps.Token},
	}
	nginx := nginxConfig{}
	endpoints := make(model.Endpoints, 0, len(targets))
	for index, target := range targets {
		exposedPort, ret := GetAvailableFrpsPort(frps.Host, portRange, target.Protocol)
		if !ret.OK {
			return frpSidecars{}, ret
		}
		// TCP 目标经 nginx 转发,frpc 的连接才能带上 proxy_protocol。
		if protocol := strings.ToLower(target.Protocol); protocol == "tcp" {
			listenPort := int32(10000 + index)
			addFrpcProxy(&frpc, protocol, "127.0.0.1", listenPort, exposedPort, true)
			nginx.Streams = append(nginx.Streams, nginxStream{
				Name:       fmt.Sprintf("upstream_%d", index),
				TargetIP:   target.IP,
				TargetPort: target.Port,
				ListenPort: listenPort,
			})
		} else {
			addFrpcProxy(&frpc, protocol, target.IP, target.Port, exposedPort, false)
		}
		endpoints = append(endpoints, model.Endpoint{
			Name:     target.Name,
			IP:       frps.Host,
			Port:     exposedPort,
			Protocol: target.Protocol,
		})
		log.Logger.Debugf(
			"Reserved frpc endpoint: %s:%d -> %s:%d protocol=%s",
			frps.Host, exposedPort, target.IP, target.Port, target.Protocol,
		)
	}
	frpcConfigData, err := frpc.String()
	if err != nil {
		return frpSidecars{}, model.FailRetVal(i18n.K8S.CreateError, "ConfigMap", configName, err)
	}
	containers := []corev1.Container{{
		Name:            FrpcContainerName,
		Resources:       sidecarResources(),
		Image:           config.Env.K8S.Frp.FrpcImage,
		ImagePullPolicy: corev1.PullIfNotPresent,
		Args:            []string{"-c", "/etc/frp/frpc.toml"},
		VolumeMounts: []corev1.VolumeMount{{
			Name:      frpcConfigVolume,
			MountPath: "/etc/frp/frpc.toml",
			SubPath:   "frpc.toml",
		}},
	}}
	volumes := []corev1.Volume{{
		Name:      frpcConfigVolume,
		ConfigMap: &corev1.ConfigMapVolumeSource{Name: configName},
	}}
	if len(nginx.Streams) > 0 {
		containers = append(containers, corev1.Container{
			Name:            NginxContainerName,
			Resources:       sidecarResources(),
			Image:           config.Env.K8S.Frp.NginxImage,
			ImagePullPolicy: corev1.PullIfNotPresent,
			VolumeMounts: []corev1.VolumeMount{{
				Name:      nginxConfigVolume,
				MountPath: "/etc/nginx/nginx.conf",
				SubPath:   "nginx.conf",
			}},
		})
		volumes = append(volumes, corev1.Volume{
			Name:      nginxConfigVolume,
			ConfigMap: &corev1.ConfigMapVolumeSource{Name: configName},
		})
	}
	return frpSidecars{
		Containers: containers,
		Volumes:    volumes,
		ConfigName: configName,
		ConfigData: map[string]string{"frpc.toml": frpcConfigData, "nginx.conf": nginx.String()},
		Endpoints:  endpoints,
	}, model.SuccessRetVal()
}

// pickFrps 随机挑一个 frps 服务器。
func pickFrps() (config.FrpsConfig, model.RetVal) {
	if len(config.Env.K8S.Frp.Frps) == 0 {
		return config.FrpsConfig{}, model.FailRetVal(i18n.K8S.NotFound, "Frps", "frps_pool", errNoFrpsServer)
	}
	index, _ := rand.Int(rand.Reader, big.NewInt(int64(len(config.Env.K8S.Frp.Frps))))
	return config.Env.K8S.Frp.Frps[index.Int64()], model.SuccessRetVal()
}

// frpsPortRange 展开 frps 允许使用的端口区间。
func frpsPortRange(frps config.FrpsConfig) []int32 {
	portRange := make([]int32, 0)
	for _, pr := range frps.Allowed {
		for i := pr.From; i <= pr.To; i++ {
			if slices.Contains(pr.Exclude, i) {
				continue
			}
			portRange = append(portRange, i)
		}
	}
	return portRange
}

// frpsPeerCIDRs 解析 frps 服务器地址,跳板的出站策略要放行它。
func frpsPeerCIDRs(host string) []string {
	if ip := net.ParseIP(host); ip != nil {
		return []string{ip.String() + "/32"}
	}
	addrs, err := net.LookupIP(host)
	if err != nil {
		log.Logger.Warningf("Failed to resolve FRPS host %s: %v", host, err)
		return nil
	}
	cidrs := make([]string, 0, len(addrs))
	for _, addr := range addrs {
		cidrs = append(cidrs, addr.String()+"/32")
	}
	return cidrs
}

func AddFrpc(ctx context.Context, victim model.Victim) (model.Victim, model.RetVal) {
	podName := fmt.Sprintf("frpc-%d", victim.ID)
	sidecars, ret := buildFrpSidecars(podName, victim.Endpoints)
	if !ret.OK {
		return victim, ret
	}
	log.Logger.Debugf("Creating frpc resources: victim_id=%d endpoints=%d", victim.ID, len(sidecars.Endpoints))
	if config.Env.K8S.CaptureEnabled {
		// The capture sidecar writes into the traffic directory, which has to
		// exist and belong to the capture user before the pod is created.
		if err := prepareCaptureTrafficDir(victim); err != nil {
			log.Logger.Warningf(
				"Failed to prepare traffic directory for victim %d: %s", victim.ID, err,
			)
		}
	}
	// The frpc pod carries its own role tag, so the victim NetworkPolicies,
	// which select on the victim role, do not match it.
	labels := VictimLabels(victim, map[string]string{RoleLabel: FrpcPodTag})
	fcm, ret := CreateConfigMap(ctx, CreateConfigMapOptions{
		Name:   sidecars.ConfigName,
		Labels: labels,
		Data:   sidecars.ConfigData,
	})
	if !ret.OK {
		return victim, ret
	}
	if fcm == nil {
		return victim, model.FailRetVal(i18n.K8S.CreateError, "ConfigMap", sidecars.ConfigName, errMissingConfigMap)
	}
	containers := sidecars.Containers
	volumes := sidecars.Volumes
	if config.Env.K8S.CaptureEnabled {
		containers = append(containers, captureContainer(victim, "frpc.pcap"))
		volumes = append(volumes, nfsVolume())
	}
	p, ret := CreatePod(ctx, CreatePodOptions{
		SubmitOnly:         true,
		Name:               podName,
		Labels:             labels,
		Containers:         containers,
		Volumes:            volumes,
		ServiceAccountName: runtimeServiceAccount(),
	})
	if !ret.OK {
		return victim, ret
	}
	if p == nil {
		return victim, model.FailRetVal(i18n.K8S.CreateError, "Pod", podName, errMissingPodObject)
	}
	victim.Resources.UIDs[podName] = string(p.UID)
	victim.Pods = append(victim.Pods, model.Pod{VictimID: victim.ID, Name: podName})
	victim.Resources.FrpcPodNames = append([]string(nil), podName)
	victim.ExposedEndpoints = sidecars.Endpoints
	log.Logger.Debugf(
		"Created frpc resources: victim_id=%d frpc_pods=%d exposed_endpoints=%d",
		victim.ID, len(victim.Resources.FrpcPodNames), len(victim.ExposedEndpoints),
	)
	return victim, model.SuccessRetVal()
}

func GetAvailableFrpsPort(host string, portRange []int32, protocol string) (int32, model.RetVal) {
	// The lock script already scans the whole range, so a failure means every
	// port is taken and narrowing the range would not help.
	port, ret := redis.LockFrpsPort(host, portRange, protocol)
	if !ret.OK {
		return 0, ret
	}
	return port, model.SuccessRetVal()
}
