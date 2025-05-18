package k8s

import (
	"CBCTF/internel/log"
	"CBCTF/internel/model"
	"CBCTF/internel/utils"
	"context"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	"strings"
)

func CreateFrpcConfig(ctx context.Context, frpsIP string, frpsPort int, token string, pod model.Pod, service *corev1.Service) (*corev1.ConfigMap, bool, string) {
	data := fmt.Sprintf("serverAddr = \"%s\"\nserverPort = %d\nauth.token = \"%s\"\n\n", frpsIP, frpsPort, token)
	tmp := make([]int32, 0)
	for _, port := range service.Spec.Ports {
		name := fmt.Sprintf("%s-%s-%d-%s", port.Protocol, pod.Name, port.Port, utils.RandStr(6))
		if utils.In(port.NodePort, tmp) {
			continue
		}
		data += fmt.Sprintf(
			"[[proxies]]\nname = \"%s\"\ntype = \"%s\"\nlocalIP = \"127.0.0.1\"\nlocalPort = %d\nremotePort = %d\n\n",
			name, strings.ToLower(string(port.Protocol)), port.Port, port.NodePort,
		)
		tmp = append(tmp, port.NodePort)
		log.Logger.Infof("Frpc started: %s:%d -> %s:%d", frpsIP, port.NodePort, pod.Name, port.Port)
	}
	return CreateConfigMap(ctx, CreateConfigMapOptions{
		PodName: pod.Name,
		Data:    map[string]string{"frpc.toml": data},
	})
}
