package k8s

import (
	"CBCTF/internel/i18n"
	"CBCTF/internel/log"
	"CBCTF/internel/model"
	"CBCTF/internel/utils"
	"context"
	"encoding/base64"
	"fmt"
	"io"
	corev1 "k8s.io/api/core/v1"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

var (
	generatorEndpoint = make(map[uint]string)
	generatorPwd      = make(map[uint]string)
	generatorMap      = make(map[string]*corev1.Pod)
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
	if pod, ok = generatorMap[generatorName]; ok && pod != nil && pod.Status.Phase == corev1.PodRunning && time.Now().Sub(pod.CreationTimestamp.Time) < 3*time.Hour {
		return pod, true, i18n.Success
	}
	if pod, ok, _ = GetPod(ctx, generatorName); pod.Status.Phase == corev1.PodRunning && time.Now().Sub(pod.CreationTimestamp.Time) < 3*time.Hour {
		log.Logger.Infof("Pod %s is already running", pod.Name)
		generatorMap[generatorName] = pod
		return pod, true, i18n.Success
	}
	generatorMap[generatorName] = nil
	if ok {
		StopGenerator(contestChallenge)
	}
	service, ok, msg := CreateService(ctx, CreateServiceOptions{
		PodName: generatorName,
		Ports:   []int32{8000},
	})
	if !ok {
		log.Logger.Warningf("Failed to create service for generator: %s", msg)
		return &corev1.Pod{}, false, msg
	}
	pwd := utils.UUID()
	pod, ok, msg = CreatePod(ctx, CreatePodOptions{
		Name: generatorName,
		Labels: map[string]string{
			"victim": generatorName,
		},
		Containers: []corev1.Container{
			{
				Name:  containerName,
				Image: contestChallenge.Challenge.GeneratorImage,
				Env: []corev1.EnvVar{
					{
						Name:  "generator_pwd",
						Value: pwd,
					},
				},
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
	generatorEndpoint[contestChallenge.ID] = fmt.Sprintf("%s:%d", pod.Status.HostIP, service.Spec.Ports[0].NodePort)
	generatorPwd[contestChallenge.ID] = pwd
	generatorMap[generatorName] = pod
	return pod, true, i18n.Success
}

// StopGenerator 停止动态附件生成器, contestChallenge 需要预加载 Challenge
func StopGenerator(contestChallenge model.ContestChallenge) (bool, string) {
	log.Logger.Infof("Stopping generator for challenge %d-%s", contestChallenge.ChallengeID, contestChallenge.Name)
	delete(generatorEndpoint, contestChallenge.ID)
	delete(generatorPwd, contestChallenge.ID)
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()
	return DeletePod(ctx, GenGeneratorName(contestChallenge.Challenge.RandID))
}

// GenerateAttachment 附加容器命令, 生成附件, model.Usage 需要预加载
func GenerateAttachment(contestChallenge model.ContestChallenge, team model.Team, teamFlagL []model.TeamFlag) (bool, string) {
	var err error
	log.Logger.Debugf("Generating attachment for team %d challenge %d", team.ID, contestChallenge.ChallengeID)
	_, ok, msg := StartGenerator(contestChallenge)
	// 附加失败则直接返回, 并尝试关闭生成器
	if !ok {
		go StopGenerator(contestChallenge)
		return false, msg
	}
	var flags string
	for _, teamFlag := range teamFlagL {
		flags += fmt.Sprintf("%s,", base64.StdEncoding.EncodeToString([]byte(teamFlag.Value)))
	}
	flags = base64.StdEncoding.EncodeToString([]byte(strings.TrimSuffix(flags, ",")))

	params := url.Values{}
	params.Add("id", fmt.Sprintf("%d", team.ID))
	params.Add("flags", flags)
	params.Add("pwd", generatorPwd[contestChallenge.ID])
	base := fmt.Sprintf("http://%s/gen?%s", generatorEndpoint[contestChallenge.ID], params.Encode())
	resp, err := http.Get(base)
	if err != nil {
		log.Logger.Warningf("Failed to generate attachment for team %d challenge %d: %v", team.ID, contestChallenge.ChallengeID, err)
		return false, i18n.ExecCommandError
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			log.Logger.Warningf("Failed to close response body: %v", err)
		}
	}(resp.Body)
	if resp.StatusCode != http.StatusOK {
		log.Logger.Warningf("Failed to generate attachment for team %d challenge %d: %s", team.ID, contestChallenge.ChallengeID, resp.Status)
		return false, i18n.ExecCommandError
	}
	file, err := os.Create(contestChallenge.Challenge.AttachmentPath(team.ID))
	if err != nil {
		log.Logger.Warningf("Failed to save attachment for team %d challenge %d: %v", team.ID, contestChallenge.ChallengeID, err)
		return false, i18n.UnknownError
	}
	defer func(file *os.File) {
		err = file.Close()
		if err != nil {
			log.Logger.Warningf("Failed to close file %s: %v", file.Name(), err)
		}
	}(file)
	if _, err = file.ReadFrom(resp.Body); err != nil {
		log.Logger.Warningf("Failed to save attachment for team %d challenge %d: %v", team.ID, contestChallenge.ChallengeID, err)
		return false, i18n.UnknownError
	}
	return true, i18n.Success
}
