package k8s

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"

	"CBCTF/internal/config"
	"CBCTF/internal/model"
)

const (
	// jumperSSHPort 是跳板容器内 sshd 监听的端口,固定不变。
	jumperSSHPort int32 = 22

	// jumperSSHPortName 是跳板端口在端点列表里的名字。
	jumperSSHPortName = "ssh"

	// jumperSSHUserEnv 与 jumperSSHKeyEnv 是跳板容器的登录账号与公钥。
	jumperSSHUserEnv = "RANGE_SSH_USER"
	jumperSSHKeyEnv  = "RANGE_SSH_KEY"
)

// jumperContainer 构造玩家跳板的 SSH 容器。
// 跳板只提供 SSH 登录,题目的环境变量、flag 与文件都不进这个容器。
func jumperContainer(access model.RangeAccess) corev1.Container {
	return corev1.Container{
		Name:            JumperContainerName,
		Image:           config.Env.K8S.Range.JumperImage,
		ImagePullPolicy: corev1.PullIfNotPresent,
		Env: []corev1.EnvVar{
			{Name: jumperSSHUserEnv, Value: access.User},
			{Name: jumperSSHKeyEnv, Value: access.PublicKey},
		},
		Resources: jumperResources(),
	}
}

// jumperResources 是跳板容器的资源,跳板要承载交互式会话,比通用侧车宽裕一些。
func jumperResources() corev1.ResourceRequirements {
	limits := corev1.ResourceList{
		corev1.ResourceCPU:    resource.MustParse("500m"),
		corev1.ResourceMemory: resource.MustParse("256Mi"),
	}
	return corev1.ResourceRequirements{
		Limits:   limits,
		Requests: workloadRequests(limits),
	}
}
