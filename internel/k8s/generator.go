package k8s

import (
	"CBCTF/internel/config"
	"CBCTF/internel/i18n"
	"CBCTF/internel/log"
	"CBCTF/internel/model"
	"CBCTF/internel/utils"
	"context"
	"encoding/base64"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	"os"
	"strings"
	"time"
)

var (
	gIPL        = make([]string, 0)
	ipGenerator = make(map[string]string)
	generatorIP = make(map[string]string)
)

func GenGeneratorName(challengeID string) string {
	return fmt.Sprintf("gen-%s-pod", challengeID)
}

// StartGenerator 启动动态附件生成器, 等待附加命令, 生成附件, model.Usage 需要预加载
func StartGenerator(usage model.Usage) (*corev1.Pod, bool, string) {
	var (
		pod           *corev1.Pod
		ok            bool
		msg           string
		err           error
		generatorName = GenGeneratorName(usage.ChallengeID)
		containerName = fmt.Sprintf("%s-%s", generatorName, strings.ToLower(utils.RandStr(5)))
	)
	if usage.Challenge.Generator == "" {
		return &corev1.Pod{}, false, "EmptyGeneratorImage"
	}
	log.Logger.Infof("Starting Generator for Challenge %s-%s", usage.ChallengeID, usage.Name)
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()
	if pod, ok, _ = GetPod(ctx, generatorName); pod.Status.Phase == corev1.PodRunning && time.Now().Sub(pod.CreationTimestamp.Time) < 3*time.Hour {
		ipGenerator[pod.Status.PodIP] = usage.ChallengeID
		generatorIP[usage.ChallengeID] = pod.Status.PodIP
		log.Logger.Infof("Pod %s is already running", pod.Name)
		return pod, true, i18n.Success
	} else {
		StopGenerator(usage)
	}
	if len(gIPL) == 0 {
		gIPL, err = utils.GetIPBlock(0, config.Env.K8S.IPPool.CIDR, config.Env.K8S.IPPool.BlockSize)
		if err != nil || len(gIPL) == 0 {
			return &corev1.Pod{}, false, i18n.EmptyIPBlock
		}
	}
	retry := 0
	ip := gIPL[retry]
	for {
		retry++
		if retry > len(gIPL)-1 {
			return &corev1.Pod{}, false, i18n.NoAvailableIP
		}
		if _, ok := ipGenerator[ip]; ok {
			ip = gIPL[retry]
			continue
		}
		ipGenerator[ip] = usage.ChallengeID
		generatorIP[usage.ChallengeID] = ip
		break
	}
	pod, ok, msg = CreatePod(ctx, CreatePodOptions{
		Name:  generatorName,
		PodIP: ip,
		Containers: []corev1.Container{
			{
				Name:    containerName,
				Image:   usage.Challenge.Generator,
				Command: []string{"sleep", "infinity"},
			},
		},
	})
	if !ok {
		return &corev1.Pod{}, false, msg
	}
	var commands []string
	if _, err = os.Stat(usage.Challenge.GeneratorPath()); err == nil {
		err = CopyToPod(generatorName, containerName, usage.Challenge.GeneratorPath(), "/root/generator.zip")
		if err != nil {
			log.Logger.Warningf("Failed to copy file: %v", err)
			return &corev1.Pod{}, false, i18n.CopyFileError
		}
		commands = append(commands, "unzip /root/generator.zip -d /root")
	} else {
		log.Logger.Warning("Generator file not found, make sure the generator docker can work correctly")
	}
	for _, command := range commands {
		log.Logger.Debugf("Executing command: %s", command)
		if _, _, err = Exec(generatorName, containerName, command, nil); err != nil {
			log.Logger.Warningf("Failed to execute command %s: %v", command, err)
			return &corev1.Pod{}, false, i18n.ExecCommandError
		}
	}
	return pod, true, i18n.Success
}

// StopGenerator 停止动态附件生成器
func StopGenerator(usage model.Usage) (bool, string) {
	log.Logger.Infof("Stopping generator for challenge %s-%s", usage.ChallengeID, usage.Name)
	if ip, ok := generatorIP[usage.ChallengeID]; ok {
		delete(generatorIP, usage.ChallengeID)
		delete(ipGenerator, ip)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()
	return DeletePod(ctx, GenGeneratorName(usage.ChallengeID))
}

// GenerateAttachment 附加容器命令, 生成附件, model.Usage 需要预加载
func GenerateAttachment(usage model.Usage, team model.Team, answer []model.Answer) (bool, string) {
	var err error
	log.Logger.Debugf("Generating attachment for team %d challenge %s", team.ID, usage.ChallengeID)
	pod, ok, msg := StartGenerator(usage)
	// 附加失败则直接返回, 并尝试关闭生成器
	if !ok {
		go StopGenerator(usage)
		return false, msg
	}
	var flags string
	for _, a := range answer {
		flags += fmt.Sprintf("%s,", base64.StdEncoding.EncodeToString([]byte(a.Value)))
	}
	flags = strings.TrimSuffix(flags, ",")
	command := fmt.Sprintf("./run.sh %d %s", team.ID, base64.StdEncoding.EncodeToString([]byte(flags)))
	log.Logger.Debugf("Executing command: %s", command)
	if _, _, err = Exec(pod.Name, pod.Spec.Containers[0].Name, command, nil); err != nil {
		log.Logger.Warningf("Failed to execute command %s: %v", command, err)
		return false, i18n.ExecCommandError
	}
	err = CopyFromPod(
		pod.Name, pod.Spec.Containers[0].Name,
		fmt.Sprintf("/root/attachments/%d.zip", team.ID),
		usage.Challenge.AttachmentPath(team.ID),
	)
	if err != nil {
		log.Logger.Warningf("Failed to copy output file: %v", err)
		return false, i18n.CopyFileError
	}
	return true, i18n.Success
}
