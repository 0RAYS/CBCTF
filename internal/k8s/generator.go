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
	"os"
)

func Generate(challenge model.Challenge) (bool, string) {
	var err error
	if challenge.Type != model.Dynamic {
		return false, "InvalidChallengeType"
	}
	if challenge.GeneratorImage == "" {
		log.Logger.Errorf("%s:%s generator image not found", challenge.ID, challenge.Name)
		return false, "EmptyGeneratorImage"
	}
	generatorPath := fmt.Sprintf("%s/challenges/%s/generator.zip", config.Env.Gin.Upload.Path, challenge.ID)
	if _, err := os.Stat(generatorPath); err != nil {
		log.Logger.Warningf("%s:%s generator.zip not found", challenge.ID, challenge.Name)
		return false, "FileNotFound"
	}
	log.Logger.Debugf("Creating pod for challenge %s:%s", challenge.Name, challenge.ID)
	podName := fmt.Sprintf("%s-generator-pod", challenge.ID)
	containerName := fmt.Sprintf("%s-generator-container", challenge.ID)
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
			RestartPolicy: corev1.RestartPolicyNever,
		},
	}
	pod, err = Client.CoreV1().Pods(NamespaceName).Create(context.TODO(), pod, metav1.CreateOptions{})
	if err != nil {
		log.Logger.Errorf("Failed to create pod: %v", err)
		return false, "CreatePodError"
	}
	defer func(pods v1.PodInterface, ctx context.Context, name string, opts metav1.DeleteOptions) {
		err = pods.Delete(ctx, name, opts)
		if err != nil {
			log.Logger.Errorf("Failed to delete pod: %v", err)
		}
		log.Logger.Debugf("%s:%s deleted successfully", challenge.Name, pod.Name)
	}(Client.CoreV1().Pods(NamespaceName), context.TODO(), pod.Name, metav1.DeleteOptions{})

	for {
		pod, err = Client.CoreV1().Pods(NamespaceName).Get(context.TODO(), pod.Name, metav1.GetOptions{})
		if err != nil {
			log.Logger.Errorf("Failed to get pod: %v", err)
			return false, "GetPodError"
		}
		if pod.Status.Phase == corev1.PodRunning {
			break
		}
		if pod.Status.Phase != corev1.PodPending {
			log.Logger.Errorf("Pod %s is not running", pod.Name)
			return false, "PodNotRunning"
		}
	}
	err = CopyToPod(NamespaceName, pod.Name, containerName, generatorPath, "/root/generator.zip")
	if err != nil {
		log.Logger.Errorf("Failed to copy file: %v", err)
		return false, "CopyFileError"
	}
	commands := []string{
		"unzip /root/generator.zip -d /root",
		"pip install -r requirements.txt",
		"python /root/generator.py 1 this_is_flag",
		"zip -r /root/attachments.zip ./attachments/*",
	}
	for _, command := range commands {
		log.Logger.Debugf("Executing command: %s", command)
		var buf bytes.Buffer
		if ExecInPod(NamespaceName, pod.Name, containerName, command, nil, &buf, nil) != nil {
			log.Logger.Errorf("Failed to execute command: %v", err)
			return false, "ExecCommandError"
		}
	}
	err = CopyFromPod(
		NamespaceName, pod.Name, containerName,
		"/root/attachments.zip",
		fmt.Sprintf("%s/attachments/%s/attachments.zip", config.Env.Gin.Upload.Path, challenge.ID),
	)
	if err != nil {
		log.Logger.Errorf("Failed to copy output file: %v", err)
		return false, "CopyFileError"
	}
	return true, ""
}
