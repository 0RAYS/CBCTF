package k8s

import (
	"CBCTF/internal/config"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"bytes"
	"context"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/utils/ptr"
	"os"
)

func GenerateAttachment(challenge model.Challenge, flag model.Flag) (bool, string) {
	var err error
	if challenge.Type != model.Dynamic {
		return false, "InvalidChallengeType"
	}
	if challenge.GeneratorImage == "" {
		return false, "EmptyGeneratorImage"
	}
	log.Logger.Debugf("Creating pod for challenge %s:%s", challenge.Name, challenge.ID)
	podName := fmt.Sprintf("%s-%d-generator-pod", challenge.ID, flag.TeamID)
	containerName := fmt.Sprintf("%s-%d-generator", challenge.ID, flag.TeamID)
	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      podName,
			Namespace: NamespaceName,
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:    containerName,
					Image:   challenge.GeneratorImage,
					Command: []string{"sleep", "infinity"},
				},
			},
			TerminationGracePeriodSeconds: ptr.To[int64](3),
			RestartPolicy:                 corev1.RestartPolicyNever,
		},
	}
	pod, err = Client.CoreV1().Pods(NamespaceName).Create(context.TODO(), pod, metav1.CreateOptions{})
	if err != nil {
		log.Logger.Warningf("Failed to create pod: %v", err)
		return false, "CreatePodError"
	}
	defer func(pods v1.PodInterface, ctx context.Context, name string, opts metav1.DeleteOptions) {
		err = pods.Delete(ctx, name, opts)
		if err != nil {
			log.Logger.Warningf("Failed to delete pod: %v", err)
		} else {
			log.Logger.Debugf("%s:%s deleted successfully", challenge.Name, pod.Name)
		}
	}(Client.CoreV1().Pods(NamespaceName), context.TODO(), pod.Name, metav1.DeleteOptions{})

	for {
		pod, err = Client.CoreV1().Pods(NamespaceName).Get(context.TODO(), pod.Name, metav1.GetOptions{})
		if err != nil {
			log.Logger.Warningf("Failed to get pod: %v", err)
			return false, "GetPodError"
		}
		if pod.Status.Phase == corev1.PodRunning {
			break
		}
		if pod.Status.Phase != corev1.PodPending {
			log.Logger.Warningf("Pod %s is not running", pod.Name)
			return false, "PodNotRunning"
		}
	}
	var commands []string
	generatorPath := challenge.GeneratorPath()
	if _, err := os.Stat(generatorPath); err == nil {
		err = CopyToPod(pod.Name, containerName, generatorPath, "/root/generator.zip")
		if err != nil {
			log.Logger.Warningf("Failed to copy file: %v", err)
			return false, "CopyFileError"
		}
		commands = append(commands, "unzip /root/generator.zip -d /root")
	} else {
		log.Logger.Warning("Generator file not found, make sure the generator docker can work correctly")
	}
	// TODO 有 RCE 的风险，虽然是在容器内
	commands = append(commands, fmt.Sprintf("python generator.py %d '%s'", flag.TeamID, flag.Value))
	for _, command := range commands {
		log.Logger.Debugf("Executing command: %s", command)
		var buf bytes.Buffer
		if ExecInPod(pod.Name, containerName, command, nil, &buf, nil) != nil {
			log.Logger.Warningf("Failed to execute command %s: %v", command, err)
			return false, "ExecCommandError"
		}
	}
	err = CopyFromPod(
		pod.Name, containerName,
		fmt.Sprintf("/root/attachments/%d.zip", flag.TeamID),
		fmt.Sprintf("%s/attachments/%s/%d.zip", config.Env.Gin.Upload.Path, challenge.ID, flag.TeamID),
	)
	if err != nil {
		log.Logger.Warningf("Failed to copy output file: %v", err)
		return false, "CopyFileError"
	}
	return true, ""
}
