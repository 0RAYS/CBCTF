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

func AddFrpc(ctx context.Context, victim model.Victim) (model.Victim, model.RetVal) {
	if len(config.Env.K8S.Frp.Frps) == 0 {
		return victim, model.FailRetVal(i18n.K8S.NotFound, "Frps", "frps_pool", errNoFrpsServer)
	}
	idxBig, _ := rand.Int(rand.Reader, big.NewInt(int64(len(config.Env.K8S.Frp.Frps))))
	frps := config.Env.K8S.Frp.Frps[idxBig.Int64()]
	log.Logger.Debugf("Creating frpc resources: victim_id=%d frps=%s endpoints=%d", victim.ID, frps.Host, len(victim.Endpoints))
	if config.Env.K8S.CaptureEnabled {
		// The capture sidecar writes into the traffic directory, which has to
		// exist and belong to the capture user before the pod is created.
		if err := prepareCaptureTrafficDir(victim); err != nil {
			log.Logger.Warningf(
				"Failed to prepare traffic directory for victim %d: %s", victim.ID, err,
			)
		}
	}
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
	podName := fmt.Sprintf("frpc-%d", victim.ID)
	frpc := frpcConfig{
		ServerAddr: frps.Host,
		ServerPort: frps.Port,
		Auth:       frpcAuth{Token: frps.Token},
	}
	nginx := nginxConfig{}
	victim.ExposedEndpoints = nil
	for index, endpoint := range victim.Endpoints {
		exposedPort, ret := GetAvailableFrpsPort(frps.Host, portRange, endpoint.Protocol)
		if !ret.OK {
			return victim, ret
		}
		// TCP endpoints go through nginx so the frpc connection can carry
		// proxy_protocol.
		if protocol := strings.ToLower(endpoint.Protocol); protocol == "tcp" {
			listenPort := int32(10000 + index)
			addFrpcProxy(&frpc, protocol, "127.0.0.1", listenPort, exposedPort, true)
			nginx.Streams = append(nginx.Streams, nginxStream{
				Name:       fmt.Sprintf("upstream_%d", index),
				TargetIP:   endpoint.IP,
				TargetPort: endpoint.Port,
				ListenPort: listenPort,
			})
		} else {
			addFrpcProxy(&frpc, protocol, endpoint.IP, endpoint.Port, exposedPort, false)
		}
		newEndpoints = append(newEndpoints, model.Endpoint{
			Name:     endpoint.Name,
			IP:       frps.Host,
			Port:     exposedPort,
			Protocol: endpoint.Protocol,
		})
		victim.ExposedEndpoints = newEndpoints
		log.Logger.Debugf(
			"Reserved frpc endpoint: victim_id=%d %s:%d -> %s:%d protocol=%s",
			victim.ID, frps.Host, exposedPort, endpoint.IP, endpoint.Port, endpoint.Protocol,
		)
	}
	frpcConfigData, err := frpc.String()
	if err != nil {
		return victim, model.FailRetVal(i18n.K8S.CreateError, "ConfigMap", podName, err)
	}
	podFrpcConfigMap[podName] = frpcConfigData
	podNginxConfigMap[podName] = nginx.String()
	frpcPodNameL = append(frpcPodNameL, podName)
	victim.ExposedEndpoints = newEndpoints
	for _, name := range frpcPodNameL {
		victim.Pods = append(victim.Pods, model.Pod{
			VictimID: victim.ID,
			Name:     name,
		})
	}
	victim.Resources.FrpcPodNames = append([]string(nil), frpcPodNameL...)
	// The frpc pod carries its own role tag, so the victim NetworkPolicies,
	// which select on the victim role, do not match it.
	labels := VictimLabels(victim, map[string]string{RoleLabel: FrpcPodTag})
	wg := utils.NewGroup(ctx)
	groupCtx := wg.Context()
	for _, name := range frpcPodNameL {
		wg.Go(func() error {
			ctx := groupCtx
			fcm, ret := CreateConfigMap(ctx, CreateConfigMapOptions{
				Name:   name,
				Labels: labels,
				Data:   map[string]string{"frpc.toml": podFrpcConfigMap[name], "nginx.conf": podNginxConfigMap[name]},
			})
			if !ret.OK || fcm == nil {
				return ret
			}
			fcmVolume := corev1.Volume{
				Name: "frpc-volume",
				ConfigMap: &corev1.ConfigMapVolumeSource{
					Name: fcm.Name,
				},
			}
			ncmVolume := corev1.Volume{
				Name: "nginx-volume",
				ConfigMap: &corev1.ConfigMapVolumeSource{
					Name: fcm.Name,
				},
			}
			nfsVolume := corev1.Volume{
				Name: nfsVolumeName,
				PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
					ClaimName: nfsVolumeClaim,
				},
			}
			containers := []corev1.Container{
				{
					Name:            FrpcContainerName,
					Resources:       sidecarResources(),
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
					Name:            NginxContainerName,
					Resources:       sidecarResources(),
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
			capture := captureContainer(victim, "frpc.pcap")
			volumes := []corev1.Volume{fcmVolume}
			if len(nginx.Streams) > 0 {
				volumes = append(volumes, ncmVolume)
			} else {
				containers = containers[:1]
			}
			if config.Env.K8S.CaptureEnabled {
				containers = append(containers, capture)
				volumes = append(volumes, nfsVolume)
			}
			options := CreatePodOptions{
				SubmitOnly:         true,
				Name:               name,
				Labels:             labels,
				Containers:         containers,
				Volumes:            volumes,
				ServiceAccountName: runtimeServiceAccount(),
			}
			p, ret := CreatePod(ctx, options)
			if !ret.OK {
				return ret
			}
			victim.Resources.UIDs[name] = string(p.UID)
			return nil
		})
	}
	if err = wg.Wait(); err != nil {
		return victim, model.FailRetVal(i18n.Common.UnknownError, "Frpc", podName, err)
	}
	log.Logger.Debugf(
		"Created frpc resources: victim_id=%d frpc_pods=%d exposed_endpoints=%d",
		victim.ID, len(frpcPodNameL), len(victim.ExposedEndpoints),
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
