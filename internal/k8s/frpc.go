package k8s

import (
	"CBCTF/internal/config"
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"CBCTF/internal/redis"
	"CBCTF/internal/utils"
	"bytes"
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"net"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/BurntSushi/toml"
	corev1 "k8s.io/api/core/v1"
)

type frpcConfig struct {
	ServerAddr string      `toml:"serverAddr"`
	ServerPort int32       `toml:"serverPort"`
	Auth       frpcAuth    `toml:"auth"`
	Proxies    []frpcProxy `toml:"proxies"`
}

type frpcAuth struct {
	Token string `toml:"token"`
}

type frpcProxy struct {
	Name       string         `toml:"name"`
	Type       string         `toml:"type"`
	LocalIP    string         `toml:"localIP"`
	LocalPort  int32          `toml:"localPort"`
	RemotePort int32          `toml:"remotePort"`
	Transport  *frpcTransport `toml:"transport,omitempty"`
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
		fmt.Fprintf(&buf, "    upstream %s {\n", stream.Name)
		fmt.Fprintf(&buf, "        server %s;\n", server)
		buf.WriteString("    }\n")
		buf.WriteString("    server {\n")
		fmt.Fprintf(&buf, "        listen %d proxy_protocol;\n", stream.ListenPort)
		fmt.Fprintf(&buf, "        proxy_pass %s;\n", stream.Name)
		buf.WriteString("    }\n")
	}
	buf.WriteString("}")
	return buf.String()
}

func newFrpcConfig(serverAddr string, serverPort int32, token string) frpcConfig {
	return frpcConfig{
		ServerAddr: serverAddr,
		ServerPort: serverPort,
		Auth: frpcAuth{
			Token: token,
		},
	}
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

func AddFrpc(ctx context.Context, victim model.Victim) (model.Victim, model.RetVal) {
	idxBig, _ := rand.Int(rand.Reader, big.NewInt(int64(len(config.Env.K8S.Frp.Frps))))
	frps := config.Env.K8S.Frp.Frps[idxBig.Int64()]
	log.Logger.Debugf("Creating frpc resources: victim_id=%d frps=%s endpoints=%d", victim.ID, frps.Host, len(victim.Endpoints))
	portRange := make([]int32, 0)
	for _, pr := range frps.Allowed {
		for i := pr.From; i <= pr.To; i++ {
			if slices.Contains(pr.Exclude, i) {
				continue
			}
			portRange = append(portRange, i)
		}
	}
	newEndpoints := make(model.Endpoints, 0)
	frpcPodNameL := make([]string, 0)
	podFrpcConfigMap := make(map[string]string)
	podNginxConfigMap := make(map[string]string)
	podVPCGWMap := make(map[string]string)
	podName := fmt.Sprintf("frpc-%d-%d-%s", victim.ContestChallengeID.V, victim.UserID, utils.RandHexStr(6))
	// 添加一个独立tag, 防止受 NetworkPolicy 影响
	frpcConfig := newFrpcConfig(frps.Host, frps.Port, frps.Token)
	nginxConfig := nginxConfig{}
	for _, endpoint := range victim.Endpoints {
		exposedPort, ret := GetAvailableFrpsPort(frps.Host, portRange, endpoint.Protocol)
		if !ret.OK {
			return victim, ret
		}
		// 对于 TCP 协议, 启用 proxy_protocol
		if protocol := strings.ToLower(endpoint.Protocol); protocol == "tcp" {
			addFrpcProxy(&frpcConfig, protocol, "127.0.0.1", endpoint.Port, exposedPort, true)
			name := utils.RandHexStr(10)
			nginxConfig.Streams = append(nginxConfig.Streams, nginxStream{Name: name, TargetIP: endpoint.IP, TargetPort: endpoint.Port, ListenPort: endpoint.Port})
		} else {
			addFrpcProxy(&frpcConfig, protocol, endpoint.IP, endpoint.Port, exposedPort, false)
		}
		newEndpoints = append(newEndpoints, model.Endpoint{
			Name:     endpoint.Name,
			IP:       frps.Host,
			Port:     exposedPort,
			Protocol: endpoint.Protocol,
		})
		log.Logger.Debugf("Reserved frpc endpoint: victim_id=%d %s:%d -> %s:%d protocol=%s", victim.ID, frps.Host, exposedPort, endpoint.IP, endpoint.Port, endpoint.Protocol)
	}
	frpcConfigData, err := frpcConfig.String()
	if err != nil {
		return victim, model.RetVal{Msg: i18n.K8S.CreateError, Attr: map[string]any{"Model": "ConfigMap", "Error": err.Error()}}
	}
	podFrpcConfigMap[podName] = frpcConfigData
	podNginxConfigMap[podName] = nginxConfig.String()
	frpcPodNameL = append(frpcPodNameL, podName)
	victim.ExposedEndpoints = newEndpoints
	victim.Pods = append(victim.Pods, func(victimID uint, podNames []string) []model.Pod {
		frpcPods := make([]model.Pod, 0, len(podNames))
		for _, podName := range podNames {
			frpcPods = append(frpcPods, model.Pod{
				VictimID: victimID,
				Name:     podName,
			})
		}
		return frpcPods
	}(victim.ID, frpcPodNameL)...)
	victim.Resources.FrpcPodNames = append(model.StringList(nil), frpcPodNameL...)
	labels := VictimLabels(victim, map[string]string{FrpcPodTag: FrpcPodTag})
	wg := utils.NewGroup(ctx)
	for _, podName := range frpcPodNameL {
		wg.Go(func() error {
			fcm, ret := CreateConfigMap(ctx, CreateConfigMapOptions{
				Name:   fmt.Sprintf("frpc-%d-%d-%s", victim.ContestChallengeID.V, victim.UserID, utils.RandHexStr(6)),
				Labels: labels,
				Data:   map[string]string{"frpc.toml": podFrpcConfigMap[podName]},
			})
			if !ret.OK || fcm == nil {
				if err, ok := ret.Attr["Error"].(string); ok {
					return fmt.Errorf("%s", err)
				}
				return fmt.Errorf("create frpc configmap failed: %s", ret.Msg)
			}
			ncm, ret := CreateConfigMap(ctx, CreateConfigMapOptions{
				Name:   fmt.Sprintf("nginx-%d-%d-%s", victim.ContestChallengeID.V, victim.UserID, utils.RandHexStr(6)),
				Labels: labels,
				Data:   map[string]string{"nginx.conf": podNginxConfigMap[podName]},
			})
			if !ret.OK || ncm == nil {
				if err, ok := ret.Attr["Error"].(string); ok {
					return fmt.Errorf("%s", err)
				}
				return fmt.Errorf("create nginx configmap failed: %s", ret.Msg)
			}
			fcmVolume := corev1.Volume{
				Name: "frpc-volume",
				VolumeSource: corev1.VolumeSource{
					ConfigMap: &corev1.ConfigMapVolumeSource{
						LocalObjectReference: corev1.LocalObjectReference{
							Name: fcm.Name,
						},
					},
				},
			}
			ncmVolume := corev1.Volume{
				Name: "nginx-volume",
				VolumeSource: corev1.VolumeSource{
					ConfigMap: &corev1.ConfigMapVolumeSource{
						LocalObjectReference: corev1.LocalObjectReference{
							Name: ncm.Name,
						},
					},
				},
			}
			nfsVolume := corev1.Volume{
				Name: nfsVolumeName,
				VolumeSource: corev1.VolumeSource{
					PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
						ClaimName: nfsVolumeName,
					},
				},
			}
			containers := []corev1.Container{
				{
					Name:            "frpc",
					Image:           config.Env.K8S.Frp.FrpcImage,
					ImagePullPolicy: corev1.PullIfNotPresent,
					Args:            []string{"-c", "/etc/frp/frpc.toml"},
					VolumeMounts: []corev1.VolumeMount{
						{
							Name:      fcmVolume.Name,
							MountPath: "/etc/frp/frpc.toml",
							SubPath:   "frpc.toml",
						},
					},
				},
				{
					Name:            "nginx",
					Image:           config.Env.K8S.Frp.NginxImage,
					ImagePullPolicy: corev1.PullIfNotPresent,
					VolumeMounts: []corev1.VolumeMount{
						{
							Name:      ncmVolume.Name,
							MountPath: "/etc/nginx/nginx.conf",
							SubPath:   "nginx.conf",
						},
					},
				},
			}
			capture := corev1.Container{
				Name:            "capture",
				Image:           config.Env.K8S.CaptureImage,
				ImagePullPolicy: corev1.PullIfNotPresent,
				Command:         []string{"/bin/sh", "-c"},
				VolumeMounts: []corev1.VolumeMount{
					{
						Name:      nfsVolumeName,
						MountPath: "/root/mnt",
						SubPath: strings.TrimPrefix(
							strings.TrimPrefix(victim.TrafficBasePath(), config.Env.Path), "/",
						),
					},
				},
				SecurityContext: &corev1.SecurityContext{
					Capabilities: &corev1.Capabilities{
						Add: []corev1.Capability{"NET_RAW", "SYS_ADMIN"},
					},
				},
				Stdin: true,
				TTY:   true,
			}
			command := "rustnet -i any --pcap-export /root/mnt/frpc.pcap"
			if _, err := os.Stat(filepath.Join(config.Env.Path, "GeoLite2-City.mmdb")); err == nil {
				command = "rustnet -i any --geoip-city /root/GeoLite2-City.mmdb  --pcap-export /root/mnt/frpc.pcap"
				capture.VolumeMounts = append(capture.VolumeMounts, corev1.VolumeMount{
					Name:      nfsVolumeName,
					MountPath: "/root/GeoLite2-City.mmdb",
					SubPath:   "GeoLite2-City.mmdb",
					ReadOnly:  true,
				})
			}
			capture.Command = append(capture.Command, command)
			containers = append(containers, capture)
			options := CreatePodOptions{
				Name:       podName,
				Labels:     labels,
				Containers: containers,
				Volumes:    []corev1.Volume{fcmVolume, ncmVolume, nfsVolume},
			}
			if gw, exists := podVPCGWMap[podName]; exists {
				// frpc pod 需要与 子网 eip 进行通信, 不能与 VPCNatGW pod 位于同一个节点, 并且跨 kube-system 与本 namespace
				options.AntiNatGWName = gw
			}
			if _, ret = CreatePod(ctx, options); !ret.OK {
				if err, ok := ret.Attr["Error"].(string); ok {
					return fmt.Errorf("%s", err)
				}
				return fmt.Errorf("create frpc pod failed: %s", ret.Msg)
			}
			return nil
		})
	}
	if err := wg.Wait(); err != nil {
		return victim, model.RetVal{Msg: i18n.Common.UnknownError, Attr: map[string]any{"Error": err.Error()}}
	}
	log.Logger.Debugf("Created frpc resources: victim_id=%d frpc_pods=%d exposed_endpoints=%d", victim.ID, len(frpcPodNameL), len(victim.ExposedEndpoints))
	return victim, model.SuccessRetVal()
}

func GetAvailableFrpsPort(host string, portRange []int32, protocol string) (int32, model.RetVal) {
	port, ret := redis.LockFrpsPort(host, portRange, protocol)
	if !ret.OK {
		if err, ok := ret.Attr["Error"]; ok && err == nil {
			portRange = slices.DeleteFunc(portRange, func(i int32) bool {
				return i == port
			})
			return GetAvailableFrpsPort(host, portRange, protocol)
		}
		return 0, ret
	}
	return port, model.SuccessRetVal()
}
