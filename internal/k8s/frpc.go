package k8s

import (
	"CBCTF/internal/config"
	"CBCTF/internal/i18n"
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
	"time"

	corev1 "k8s.io/api/core/v1"
)

type CreateFrpcPodResult struct {
	Name string
	OK   bool
	MSG  string
}

func CreateFrpc(victim model.Victim) (model.Endpoints, string, bool, string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	idxBig, _ := rand.Int(rand.Reader, big.NewInt(int64(len(config.Env.K8S.Frpc.Frps))))
	frps := config.Env.K8S.Frpc.Frps[idxBig.Int64()]
	portRange := make([]int32, 0)
	for _, pr := range frps.AllowedPorts {
		for i := pr.From; i <= pr.To; i++ {
			if slices.Contains(pr.Exclude, i) {
				continue
			}
			portRange = append(portRange, i)
		}
	}
	newEndpoints := make(model.Endpoints, 0)
	podName := fmt.Sprintf("frpc-%s", utils.RandStr(20))
	// 添加一个独立tag, 防止受 NetworkPolicy 影响
	labels := map[string]string{
		"victim_id":            strconv.Itoa(int(victim.ID)),
		"user_id":              strconv.Itoa(int(victim.UserID.V)),
		"team_id":              strconv.Itoa(int(victim.TeamID.V)),
		"challenge_id":         strconv.Itoa(int(victim.ChallengeID)),
		"contest_challenge_id": strconv.Itoa(int(victim.ContestChallengeID.V)),
		FrpcPodTag:             podName,
	}
	data := fmt.Sprintf("serverAddr = \"%s\"\nserverPort = %d\nauth.token = \"%s\"\n\n", frps.Host, frps.Port, frps.Token)
	for _, endpoint := range victim.Endpoints {
		exposedPort, ok, msg := GetAvailableFrpsPort(frps.Host, portRange, endpoint.Protocol)
		if !ok {
			return nil, "", false, msg
		}
		data += fmt.Sprintf(
			"[[proxies]]\nname = \"%s\"\ntype = \"%s\"\nlocalIP = \"%s\"\nlocalPort = %d\nremotePort = %d\n\n",
			utils.RandStr(10), strings.ToLower(endpoint.Protocol), endpoint.IP, endpoint.Port, exposedPort,
		)
		newEndpoints = append(newEndpoints, model.Endpoint{
			IP:       frps.Host,
			Port:     exposedPort,
			Protocol: endpoint.Protocol,
		})
		log.Logger.Infof("Frpc started: %s:%d -> %s:%d", frps.Host, exposedPort, endpoint.IP, endpoint.Port)
	}
	antiAffinity := make([]string, 0)
	for _, subnet := range victim.VPC.Subnets {
		if subnet.NatGateway == nil {
			continue
		}
		antiAffinity = append(antiAffinity, fmt.Sprintf("vpc-nat-gw-%s", subnet.NatGateway.Name))
	}
	cm, ok, msg := CreateConfigMap(ctx, CreateConfigMapOptions{
		Name:   fmt.Sprintf("cm-%s", utils.RandStr(20)),
		Labels: labels,
		Data:   map[string]string{"frpc.toml": data},
	})
	if !ok {
		return nil, "", false, msg
	}
	cmVolume := corev1.Volume{
		Name: fmt.Sprintf("vol-%s", utils.RandStr(20)),
		VolumeSource: corev1.VolumeSource{
			ConfigMap: &corev1.ConfigMapVolumeSource{
				LocalObjectReference: corev1.LocalObjectReference{
					Name: cm.Name,
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
			Name:  fmt.Sprintf("frpc-%s", utils.RandStr(20)),
			Image: config.Env.K8S.Frpc.Image,
			Args:  []string{"-c", "/etc/frp/frpc.toml"},
			VolumeMounts: []corev1.VolumeMount{
				{
					Name:      cmVolume.Name,
					MountPath: "/etc/frp/frpc.toml",
					SubPath:   "frpc.toml",
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
		Volumes:    []corev1.Volume{cmVolume, nfsVolume},
	}
	if len(antiAffinity) > 0 {
		options.PodAntiAffinity = map[string][]string{"app": antiAffinity}
	}
	_, ok, msg = CreatePod(ctx, options)
	log.Logger.Debugf("Create Pod %s: %s", podName, msg)
	if !ok {
		return nil, "", false, msg
	}
	return newEndpoints, podName, true, i18n.Success
}

func GetAvailableFrpsPort(host string, portRange []int32, protocol string) (int32, bool, string) {
	idxBig, _ := rand.Int(rand.Reader, big.NewInt(int64(len(portRange))))
	port := portRange[idxBig.Int64()]
	locked, err := redis.IsFrpsPortLocked(host, port, protocol)
	if err != nil {
		log.Logger.Warningf("Failed to check if port %d is locked: %s", port, err)
		return 0, false, i18n.RedisError
	}
	if locked {
		portRange = slices.DeleteFunc(portRange, func(i int32) bool {
			return i == port
		})
		return GetAvailableFrpsPort(host, portRange, protocol)
	}
	if err = redis.LockFrpsPort(host, port, protocol); err != nil {
		return 0, false, i18n.RedisError
	}
	return port, true, i18n.Success
}
