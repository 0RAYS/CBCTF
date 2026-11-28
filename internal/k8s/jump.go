package k8s

import (
	"strconv"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"

	"CBCTF/internal/config"
)

const (
	// jumpContainerName 是跳板 Pod 里 SSH 容器的名字。
	jumpContainerName = "ssh"

	// jumpSSHVolume 是跳板 SSH 配置的卷名。
	jumpSSHVolume = "jump-ssh"

	// jumpSSHMountPath 是 SSH 配置的挂载点,镜像里的入口脚本按这个路径读取。
	jumpSSHMountPath = "/run/range/ssh"

	// jumpSSHPortEnv 告诉镜像里的 sshd 监听哪个端口。
	jumpSSHPortEnv = "RANGE_SSH_PORT"
)

// jumpContainer 构造玩家跳板的 SSH 容器。
// 跳板只提供 SSH 登录,题目的环境变量、flag 与文件都不进这个容器。
func jumpContainer(sshPort int32) corev1.Container {
	return corev1.Container{
		Name:            jumpContainerName,
		Image:           config.Env.K8S.Range.JumpImage,
		ImagePullPolicy: corev1.PullIfNotPresent,
		Env:             []corev1.EnvVar{{Name: jumpSSHPortEnv, Value: strconv.Itoa(int(sshPort))}},
		Resources:       jumpResources(),
		VolumeMounts: []corev1.VolumeMount{{
			Name:      jumpSSHVolume,
			MountPath: jumpSSHMountPath,
			ReadOnly:  true,
		}},
	}
}

// jumpSSHConfig 是跳板 SSH 容器的配置内容。
func jumpSSHConfig(user, publicKey string) map[string]string {
	return map[string]string{
		"user":            user,
		"authorized_keys": publicKey,
	}
}

// jumpResources 是跳板容器的资源,跳板要承载交互式会话,比通用侧车宽裕一些。
func jumpResources() corev1.ResourceRequirements {
	limits := corev1.ResourceList{
		corev1.ResourceCPU:    resource.MustParse("500m"),
		corev1.ResourceMemory: resource.MustParse("256Mi"),
	}
	return corev1.ResourceRequirements{
		Limits:   limits,
		Requests: workloadRequests(limits),
	}
}
