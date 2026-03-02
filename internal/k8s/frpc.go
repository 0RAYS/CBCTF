package k8s

import (
	"CBCTF/internal/config"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"CBCTF/internal/redis"
	"CBCTF/internal/utils"
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"slices"
	"strconv"
	"strings"

	corev1 "k8s.io/api/core/v1"
)

const (
	frpcHeaderTemplate = `
serverAddr = "%s"
serverPort = %d
auth.token = "%s"
`
	frpcItemTemplate = `
[[proxies]]
name = "%s"
type = "%s"
localIP = "%s"
localPort = %d
remotePort = %d
%s
`
	proxyProtocol = `transport.proxyProtocolVersion = "v2"`
)

const (
	nginxHeaderTemplate = `
worker_processes auto;
events {
    worker_connections 1024;
}
stream {
%s
}
`
	nginxItemTemplate = `
    upstream %s {
        server %s:%d;
    }
    server {
        listen %d proxy_protocol;
        proxy_pass %s;
    }
`
)

type CreateFrpcPodResult struct {
	Name string
	OK   bool
	MSG  string
}

func CreateFrpc(ctx context.Context, victim model.Victim) (model.Endpoints, []string, model.RetVal) {
	idxBig, _ := rand.Int(rand.Reader, big.NewInt(int64(len(config.Env.K8S.Frp.Frps))))
	frps := config.Env.K8S.Frp.Frps[idxBig.Int64()]
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
	if len(victim.VPC.Subnets) == 0 {
		podName := fmt.Sprintf("frpc-%s", utils.RandStr(20))
		// 添加一个独立tag, 防止受 NetworkPolicy 影响
		frpcConfig := fmt.Sprintf(frpcHeaderTemplate, frps.Host, frps.Port, frps.Token)
		nginxConfig := ""
		for _, endpoint := range victim.Endpoints {
			exposedPort, ret := GetAvailableFrpsPort(frps.Host, portRange, endpoint.Protocol)
			if !ret.OK {
				return nil, nil, ret
			}
			// 对于 TCP 协议, 启用 proxy_protocol
			if protocol := strings.ToLower(endpoint.Protocol); protocol == "tcp" {
				frpcConfig += fmt.Sprintf(
					frpcItemTemplate,
					utils.RandStr(10), strings.ToLower(endpoint.Protocol), "127.0.0.1", endpoint.Port, exposedPort, proxyProtocol,
				)
				name := utils.RandStr(10)
				nginxConfig += fmt.Sprintf(
					nginxItemTemplate,
					name, endpoint.IP, endpoint.Port, endpoint.Port, name,
				)
			} else {
				frpcConfig += fmt.Sprintf(
					frpcItemTemplate,
					utils.RandStr(10), strings.ToLower(endpoint.Protocol), endpoint.IP, endpoint.Port, exposedPort, "",
				)
			}
			newEndpoints = append(newEndpoints, model.Endpoint{
				IP:       frps.Host,
				Port:     exposedPort,
				Protocol: endpoint.Protocol,
			})
			log.Logger.Infof("Frpc started: %s:%d -> %s:%d", frps.Host, exposedPort, endpoint.IP, endpoint.Port)
		}
		podFrpcConfigMap[podName] = frpcConfig
		podNginxConfigMap[podName] = fmt.Sprintf(nginxHeaderTemplate, nginxConfig)
		frpcPodNameL = append(frpcPodNameL, podName)
	} else {
		for _, subnet := range victim.VPC.Subnets {
			if subnet.NatGateway == nil {
				continue
			}
			needFrpc := false
			podName := fmt.Sprintf("frpc-%s", utils.RandStr(20))
			frpcConfig := fmt.Sprintf(frpcHeaderTemplate, frps.Host, frps.Port, frps.Token)
			nginxConfig := ""
			for _, eip := range subnet.NatGateway.EIPs {
				for _, dnat := range eip.DNats {
					exposedPort, ret := GetAvailableFrpsPort(frps.Host, portRange, dnat.Protocol)
					if !ret.OK {
						return nil, nil, ret
					}
					// 对于 TCP 协议, 启用 proxy_protocol
					if protocol := strings.ToLower(dnat.Protocol); protocol == "tcp" {
						frpcConfig += fmt.Sprintf(
							frpcItemTemplate,
							utils.RandStr(10), strings.ToLower(dnat.Protocol), "127.0.0.1", dnat.ExternalPort, exposedPort, proxyProtocol,
						)
						name := utils.RandStr(10)
						nginxConfig += fmt.Sprintf(
							nginxItemTemplate,
							name, eip.IP, dnat.ExternalPort, dnat.ExternalPort, name,
						)
					} else {
						frpcConfig += fmt.Sprintf(
							frpcItemTemplate,
							utils.RandStr(10), strings.ToLower(dnat.Protocol), eip.IP, dnat.ExternalPort, exposedPort, "",
						)
					}
					newEndpoints = append(newEndpoints, model.Endpoint{
						IP:       frps.Host,
						Port:     exposedPort,
						Protocol: dnat.Protocol,
					})
					log.Logger.Infof("Frpc started: %s:%d -> %s:%d", frps.Host, exposedPort, eip.IP, dnat.ExternalPort)
					needFrpc = true
				}
			}
			if !needFrpc {
				continue
			}
			podFrpcConfigMap[podName] = frpcConfig
			podNginxConfigMap[podName] = fmt.Sprintf(nginxHeaderTemplate, nginxConfig)
			podVPCGWMap[podName] = subnet.NatGateway.Name
			frpcPodNameL = append(frpcPodNameL, podName)
		}
	}
	labels := map[string]string{
		"victim_id":            strconv.Itoa(int(victim.ID)),
		"user_id":              strconv.Itoa(int(victim.UserID)),
		"team_id":              strconv.Itoa(int(victim.TeamID.V)),
		"contest_id":           strconv.Itoa(int(victim.ContestID.V)),
		"challenge_id":         strconv.Itoa(int(victim.ChallengeID)),
		"contest_challenge_id": strconv.Itoa(int(victim.ContestChallengeID.V)),
		FrpcPodTag:             FrpcPodTag,
	}
	for _, podName := range frpcPodNameL {
		fcm, ret := CreateConfigMap(ctx, CreateConfigMapOptions{
			Name:   fmt.Sprintf("cm-%s", utils.RandStr(20)),
			Labels: labels,
			Data:   map[string]string{"frpc.toml": podFrpcConfigMap[podName]},
		})
		if !ret.OK || fcm == nil {
			return nil, nil, ret
		}
		ncm, ret := CreateConfigMap(ctx, CreateConfigMapOptions{
			Name:   fmt.Sprintf("cm-%s", utils.RandStr(20)),
			Labels: labels,
			Data:   map[string]string{"nginx.conf": podNginxConfigMap[podName]},
		})
		if !ret.OK || ncm == nil {
			return nil, nil, ret
		}
		fcmVolume := corev1.Volume{
			Name: fmt.Sprintf("vol-%s", utils.RandStr(20)),
			VolumeSource: corev1.VolumeSource{
				ConfigMap: &corev1.ConfigMapVolumeSource{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: fcm.Name,
					},
				},
			},
		}
		ncmVolume := corev1.Volume{
			Name: fmt.Sprintf("vol-%s", utils.RandStr(20)),
			VolumeSource: corev1.VolumeSource{
				ConfigMap: &corev1.ConfigMapVolumeSource{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: ncm.Name,
					},
				},
			},
		}
		nfsVolume := corev1.Volume{
			Name: fmt.Sprintf("vol-%s", utils.RandStr(20)),
			VolumeSource: corev1.VolumeSource{
				PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
					ClaimName: nfsVolumeName,
				},
			},
		}
		containers := []corev1.Container{
			{
				Name:  "frpc",
				Image: config.Env.K8S.Frp.FrpcImage,
				Args:  []string{"-c", "/etc/frp/frpc.toml"},
				VolumeMounts: []corev1.VolumeMount{
					{
						Name:      fcmVolume.Name,
						MountPath: "/etc/frp/frpc.toml",
						SubPath:   "frpc.toml",
					},
				},
			},
			{
				Name:  "nginx",
				Image: config.Env.K8S.Frp.NginxImage,
				VolumeMounts: []corev1.VolumeMount{
					{
						Name:      ncmVolume.Name,
						MountPath: "/etc/nginx/nginx.conf",
						SubPath:   "nginx.conf",
					},
				},
			},
			{
				Name:    "tcpdump",
				Image:   config.Env.K8S.TCPDumpImage,
				Command: []string{"/bin/sh", "-c", "tcpdump -i any -w /root/mnt/frpc.pcap"},
				VolumeMounts: []corev1.VolumeMount{
					{
						Name:      nfsVolume.Name,
						MountPath: "/root/mnt",
						SubPath: strings.TrimPrefix(
							strings.TrimPrefix(victim.TrafficBasePath(), config.Env.Path), "/",
						),
					},
				},
			},
		}
		options := CreatePodOptions{
			Name:       podName,
			Labels:     labels,
			Containers: containers,
			Volumes:    []corev1.Volume{fcmVolume, ncmVolume, nfsVolume},
		}
		if gw, exists := podVPCGWMap[podName]; exists {
			options.PodAntiAffinity = map[string]string{"app": fmt.Sprintf("vpc-nat-gw-%s", gw)}
		}
		if _, ret = CreatePod(ctx, options); !ret.OK {
			return nil, nil, ret
		}
	}
	return newEndpoints, frpcPodNameL, model.SuccessRetVal()
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
