package k8s

import (
	"CBCTF/internal/config"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"
	"os"
	"time"
)

func StartGenerator(challenge model.Challenge) (*corev1.Pod, bool, string) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()
	log.Logger.Infof("Starting generator for challenge %s-%s", challenge.ID, challenge.Name)
	var err error
	if challenge.Type != model.Dynamic {
		return &corev1.Pod{}, false, "InvalidChallengeType"
	}
	if challenge.GeneratorImage == "" {
		return &corev1.Pod{}, false, "EmptyGeneratorImage"
	}
	log.Logger.Debugf("Creating pod for challenge %s:%s", challenge.Name, challenge.ID)
	podName := fmt.Sprintf("generator-%s-pod", challenge.ID)
	pod, ok, _ := GetPod(podName)
	if ok && pod.Status.Phase == corev1.PodRunning {
		log.Logger.Infof("Pod %s is already running", pod.Name)
		return pod, true, "Success"
	}
	containerName := fmt.Sprintf("generator-%s", challenge.ID)
	pod = &corev1.Pod{
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
	pod, err = Client.CoreV1().Pods(NamespaceName).Create(ctx, pod, metav1.CreateOptions{})
	if err != nil {
		log.Logger.Warningf("Failed to create pod: %v", err)
		return &corev1.Pod{}, false, "CreatePodError"
	}
	for {
		pod, ok, _ = GetPod(podName)
		if !ok {
			log.Logger.Warningf("Failed to get pod: %v", err)
			return &corev1.Pod{}, false, "GetPodError"
		}
		if pod.Status.Phase == corev1.PodRunning {
			break
		}
		if pod.Status.Phase != corev1.PodPending {
			log.Logger.Warningf("Pod %s is not running", pod.Name)
			return &corev1.Pod{}, false, "PodNotRunning"
		}
	}
	var commands []string
	generatorPath := challenge.GeneratorPath()
	if _, err := os.Stat(generatorPath); err == nil {
		err = CopyToPod(podName, containerName, generatorPath, "/root/generator.zip")
		if err != nil {
			log.Logger.Warningf("Failed to copy file: %v", err)
			return &corev1.Pod{}, false, "CopyFileError"
		}
		commands = append(commands, "unzip /root/generator.zip -d /root")
	} else {
		log.Logger.Warning("Generator file not found, make sure the generator docker can work correctly")
	}
	for _, command := range commands {
		log.Logger.Debugf("Executing command: %s", command)
		var buf bytes.Buffer
		if ExecInPod(podName, containerName, command, nil, &buf, nil) != nil {
			log.Logger.Warningf("Failed to execute command %s: %v", command, err)
			return &corev1.Pod{}, false, "ExecCommandError"
		}
	}
	return pod, true, "Success"
}

func StopGenerator(challenge model.Challenge) (bool, string) {
	log.Logger.Infof("Stopping generator for challenge %s-%s", challenge.ID, challenge.Name)
	podName := fmt.Sprintf("generator-%s-pod", challenge.ID)
	return DeletePod(podName)
}

// GenerateAttachment 启动容器, 生成附件
func GenerateAttachment(challenge model.Challenge, flag model.Flag) (bool, string) {
	var err error
	log.Logger.Debugf("Generating attachment for team %d challenge %s", flag.TeamID, flag.ChallengeID)
	pod, ok, msg := StartGenerator(challenge)
	if !ok {
		_, _ = StopGenerator(challenge)
		return false, msg
	}
	command := fmt.Sprintf("./run.sh %d %s", flag.TeamID, base64.StdEncoding.EncodeToString([]byte(flag.Value)))
	log.Logger.Debugf("Executing command: %s", command)
	var buf bytes.Buffer
	if ExecInPod(pod.Name, pod.Spec.Containers[0].Name, command, nil, &buf, nil) != nil {
		log.Logger.Warningf("Failed to execute command %s: %v", command, err)
		return false, "ExecCommandError"
	}
	err = CopyFromPod(
		pod.Name, pod.Spec.Containers[0].Name,
		fmt.Sprintf("/root/attachments/%d.zip", flag.TeamID),
		fmt.Sprintf("%s/attachments/%s/%d.zip", config.Env.Gin.Upload.Path, challenge.ID, flag.TeamID),
	)
	if err != nil {
		log.Logger.Warningf("Failed to copy output file: %v", err)
		return false, "CopyFileError"
	}
	return true, ""
}
