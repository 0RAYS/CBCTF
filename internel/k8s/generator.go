package k8s

import (
	"CBCTF/internel/config"
	"CBCTF/internel/log"
	"CBCTF/internel/model"
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	"os"
	"time"
)

// StartGenerator 启动动态附件生成器, 等待附加命令, 生成附件, model.Usage 需要预加载
func StartGenerator(usage model.Usage) (*corev1.Pod, bool, string) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()
	log.Logger.Infof("Starting generator for challenge %s-%s", usage.ChallengeID, usage.Name)
	var err error
	if usage.Challenge.Type != model.DynamicChallenge {
		return &corev1.Pod{}, false, "InvalidChallengeType"
	}
	if usage.Challenge.Generator == "" {
		return &corev1.Pod{}, false, "EmptyGeneratorImage"
	}
	log.Logger.Debugf("Creating Pod for challenge %s:%s", usage.Name, usage.ChallengeID)
	podName := fmt.Sprintf("generator-%s-pod", usage.ChallengeID)
	// 如果已启动并健康运行, 直接返回
	pod, ok, _ := GetPod(ctx, podName)
	if ok && pod.Status.Phase == corev1.PodRunning {
		log.Logger.Infof("Pod %s is already running", pod.Name)
		return pod, true, "Success"
	}
	containerName := fmt.Sprintf("generator-%s", usage.ChallengeID)
	containers := []corev1.Container{
		{
			Name:    containerName,
			Image:   usage.Challenge.Generator,
			Command: []string{"sleep", "infinity"},
		},
	}
	pod, ok, msg := CreatePod(ctx, podName, containers)
	if !ok {
		return &corev1.Pod{}, false, msg
	}
	var commands []string
	generatorPath := usage.Challenge.GeneratorPath()
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

// StopGenerator 停止动态附件生成器
func StopGenerator(usage model.Usage) (bool, string) {
	log.Logger.Infof("Stopping generator for challenge %s-%s", usage.ChallengeID, usage.Name)
	podName := fmt.Sprintf("generator-%s-pod", usage.ChallengeID)
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()
	return DeletePod(ctx, podName)
}

// GenerateAttachment 附加容器命令, 生成附件, model.Usage 需要预加载
func GenerateAttachment(usage model.Usage, answer model.Answer) (bool, string) {
	var err error
	log.Logger.Debugf("Generating attachment for team %d challenge %s", answer.TeamID, usage.ChallengeID)
	pod, ok, msg := StartGenerator(usage)
	// 附加失败则直接返回, 并尝试关闭生成器
	if !ok {
		_, _ = StopGenerator(usage)
		return false, msg
	}
	command := fmt.Sprintf("./run.sh %d %s", answer.TeamID, base64.StdEncoding.EncodeToString([]byte(answer.Value)))
	log.Logger.Debugf("Executing command: %s", command)
	var buf bytes.Buffer
	if ExecInPod(pod.Name, pod.Spec.Containers[0].Name, command, nil, &buf, nil) != nil {
		log.Logger.Warningf("Failed to execute command %s: %v", command, err)
		return false, "ExecCommandError"
	}
	err = CopyFromPod(
		pod.Name, pod.Spec.Containers[0].Name,
		fmt.Sprintf("/root/attachments/%d.zip", answer.TeamID),
		fmt.Sprintf("%s/attachments/%s/%d.zip", config.Env.Path, usage.ChallengeID, answer.TeamID),
	)
	if err != nil {
		log.Logger.Warningf("Failed to copy output file: %v", err)
		return false, "CopyFileError"
	}
	return true, "Success"
}
