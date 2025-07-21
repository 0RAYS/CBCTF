package k8s

import (
	"CBCTF/internel/config"
	"CBCTF/internel/i18n"
	"CBCTF/internel/log"
	"CBCTF/internel/model"
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
	labels := map[string]string{
		"user_id":              fmt.Sprintf("%d", victim.UserID),
		"team_id":              fmt.Sprintf("%d", victim.TeamID),
		"contest_challenge_id": fmt.Sprintf("%d", victim.ContestChallengeID),
	}
	idxBig, _ := rand.Int(rand.Reader, big.NewInt(int64(len(config.Env.K8S.Frpc.Frps))))
	frps := config.Env.K8S.Frpc.Frps[idxBig.Int64()]
	data := fmt.Sprintf("serverAddr = \"%s\"\nserverPort = %d\nauth.token = \"%s\"\n\n", frps.Host, frps.Port, frps.Token)
	tmp := make([]int32, 0)
	for _, endpoint := range victim.Endpoints {
		if !slices.Contains(tmp, endpoint.Port) {
			data += fmt.Sprintf(
				"[[proxies]]\nname = \"%s\"\ntype = \"%s\"\nlocalIP = \"%s\"\nlocalPort = %d\nremotePort = %d\n\n",
				utils.RandStr(10), strings.ToLower(endpoint.Protocol), endpoint.IP, endpoint.Port, endpoint.Port,
			)
			tmp = append(tmp, endpoint.Port)
			log.Logger.Infof("Frpc started: %s:%d -> %s:%d", frps.Host, endpoint.Port, endpoint.IP, endpoint.Port)
		}
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
		Name:       fmt.Sprintf("pod-%s", utils.RandStr(20)),
		Labels:     labels,
		Containers: []corev1.Container{frpc},
		Volumes:    []corev1.Volume{volume},
	})
	if !ok {
		return model.Endpoints{}, false, msg
	}
	newEndpoints := make(model.Endpoints, 0)
	for _, endpoint := range victim.Endpoints {
		newEndpoints = append(newEndpoints, model.Endpoint{
			IP:       frps.Host,
			Port:     endpoint.Port,
			Protocol: endpoint.Protocol,
		})
	}
	return newEndpoints, true, i18n.Success
}
