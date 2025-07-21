package k8s

import (
	"CBCTF/internel/config"
	"CBCTF/internel/i18n"
	"CBCTF/internel/log"
	"CBCTF/internel/model"
	"CBCTF/internel/redis"
	"CBCTF/internel/utils"
	"context"
	"crypto/rand"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	"math/big"
	"slices"
	"strings"
	"time"
)

func CreateFrpc(victim model.Victim) (model.Endpoints, bool, string) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()
	podName := fmt.Sprintf("pod-%s", utils.RandStr(20))
	// 添加一个独立tag, 防止受 NetworkPolicy 影响
	labels := map[string]string{
		"user_id":              fmt.Sprintf("%d", victim.UserID),
		"team_id":              fmt.Sprintf("%d", victim.TeamID),
		"contest_challenge_id": fmt.Sprintf("%d", victim.ContestChallengeID),
		FrpcPodTag:             podName,
	}
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
	data := fmt.Sprintf("serverAddr = \"%s\"\nserverPort = %d\nauth.token = \"%s\"\n\n", frps.Host, frps.Port, frps.Token)
	newEndpoints := make(model.Endpoints, 0)
	for _, endpoint := range victim.Endpoints {
		exposedPort, ok, msg := GetAvailablePort(frps.Host, portRange, endpoint.Protocol)
		if !ok {
			return model.Endpoints{}, false, msg
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
	cm, ok, msg := CreateConfigMap(ctx, CreateConfigMapOptions{
		Name:   fmt.Sprintf("cm-%s", utils.RandStr(20)),
		Labels: labels,
		Data:   map[string]string{"frpc.toml": data},
	})
	if !ok {
		return model.Endpoints{}, false, msg
	}
	volume := corev1.Volume{
		Name: fmt.Sprintf("vol-%s", utils.RandStr(20)),
		VolumeSource: corev1.VolumeSource{
			ConfigMap: &corev1.ConfigMapVolumeSource{
				LocalObjectReference: corev1.LocalObjectReference{
					Name: cm.Name,
				},
			},
		},
	}
	frpc := corev1.Container{
		Name:  fmt.Sprintf("frpc-%s", utils.RandStr(20)),
		Image: config.Env.K8S.Frpc.Image,
		Args:  []string{"-c", "/etc/frp/frpc.toml"},
		VolumeMounts: []corev1.VolumeMount{
			{
				Name:      volume.Name,
				MountPath: "/etc/frp/frpc.toml",
				SubPath:   "frpc.toml",
			},
		},
	}
	_, ok, msg = CreatePod(ctx, CreatePodOptions{
		Name:       podName,
		Labels:     labels,
		Containers: []corev1.Container{frpc},
		Volumes:    []corev1.Volume{volume},
	})
	if !ok {
		return model.Endpoints{}, false, msg
	}
	return newEndpoints, true, i18n.Success
}

func GetAvailablePort(host string, portRange []int32, protocol string) (int32, bool, string) {
	idxBig, _ := rand.Int(rand.Reader, big.NewInt(int64(len(portRange))))
	port := portRange[idxBig.Int64()]
	locked, err := redis.IsFrpsPortLocked(host, port, protocol)
	if err != nil {
		log.Logger.Warningf("Failed to check if port %d is locked: %v", port, err)
		return 0, false, i18n.RedisError
	}
	if locked {
		portRange = slices.DeleteFunc(portRange, func(i int32) bool {
			return i == port
		})
		return GetAvailablePort(host, portRange, protocol)
	}
	if err = redis.LockFrpsPort(host, port, protocol); err != nil {
		return 0, false, i18n.RedisError
	}
	return port, true, i18n.Success
}
