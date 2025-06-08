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
	ipGenerator = make(map[string]uint)
	generatorIP = make(map[uint]string)
)

func GenGeneratorName(challengeRandID string) string {
	return fmt.Sprintf("gen-%s-pod", challengeRandID)
}

// StartGenerator 启动动态附件生成器, 等待附加命令, 生成附件, contestChallenge 需要预加载 Challenge
func StartGenerator(contestChallenge model.ContestChallenge) (*corev1.Pod, bool, string) {
	var (
		pod           *corev1.Pod
		ok            bool
		msg           string
		err           error
		generatorName = GenGeneratorName(contestChallenge.Challenge.RandID)
		containerName = fmt.Sprintf("%s-%s", generatorName, strings.ToLower(utils.RandStr(5)))
	)
	if contestChallenge.Challenge.GeneratorImage == "" {
		return &corev1.Pod{}, false, i18n.InvalidDockerImage
	}
	log.Logger.Infof("Starting Generator for Challenge %d-%s", contestChallenge.ChallengeID, contestChallenge.Name)
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()
	if pod, ok, _ = GetPod(ctx, generatorName); pod.Status.Phase == corev1.PodRunning && time.Now().Sub(pod.CreationTimestamp.Time) < 3*time.Hour {
		ipGenerator[pod.Status.PodIP] = contestChallenge.ChallengeID
		generatorIP[contestChallenge.ChallengeID] = pod.Status.PodIP
		log.Logger.Infof("Pod %s is already running", pod.Name)
		return pod, true, i18n.Success
	}
	StopGenerator(contestChallenge)
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
			return &corev1.Pod{}, false, i18n.EmptyIPBlock
		}
		if _, ok := ipGenerator[ip]; ok {
			ip = gIPL[retry]
			continue
		}
		ipGenerator[ip] = contestChallenge.ChallengeID
		generatorIP[contestChallenge.ChallengeID] = ip
		break
	}
	pod, ok, msg = CreatePod(ctx, CreatePodOptions{
		Name:  generatorName,
		PodIP: ip,
		Containers: []corev1.Container{
			{
				Name:    containerName,
				Image:   contestChallenge.Challenge.GeneratorImage,
				Command: []string{"sleep", "infinity"},
			},
		},
	})
	if !ok {
		return &corev1.Pod{}, false, msg
	}
	var commands []string
	if _, err = os.Stat(contestChallenge.Challenge.GeneratorPath()); err == nil {
		err = CopyToPod(generatorName, containerName, contestChallenge.Challenge.GeneratorPath(), "/root/generator.zip")
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

// StopGenerator 停止动态附件生成器, contestChallenge 需要预加载 Challenge
func StopGenerator(contestChallenge model.ContestChallenge) (bool, string) {
	log.Logger.Infof("Stopping generator for challenge %d-%s", contestChallenge.ChallengeID, contestChallenge.Name)
	if ip, ok := generatorIP[contestChallenge.ChallengeID]; ok {
		delete(generatorIP, contestChallenge.ChallengeID)
		delete(ipGenerator, ip)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()
	return DeletePod(ctx, GenGeneratorName(contestChallenge.Challenge.RandID))
}

// GenerateAttachment 附加容器命令, 生成附件, model.Usage 需要预加载
func GenerateAttachment(contestChallenge model.ContestChallenge, team model.Team, teamFlagL []model.TeamFlag) (bool, string) {
	var err error
	log.Logger.Debugf("Generating attachment for team %d challenge %d", team.ID, contestChallenge.ChallengeID)
	pod, ok, msg := StartGenerator(contestChallenge)
	// 附加失败则直接返回, 并尝试关闭生成器
	if !ok {
		go StopGenerator(contestChallenge)
		return false, msg
	}
	var flags string
	for _, teamFlag := range teamFlagL {
		flags += fmt.Sprintf("%s,", base64.StdEncoding.EncodeToString([]byte(teamFlag.Value)))
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
		contestChallenge.Challenge.AttachmentPath(team.ID),
	)
	if err != nil {
		log.Logger.Warningf("Failed to copy output file: %v", err)
		return false, i18n.CopyFileError
	}
	return true, i18n.Success
}
