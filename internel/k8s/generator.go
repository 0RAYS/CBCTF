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
	"math/rand"
	"os"
	"slices"
	"strings"
	"sync"
	"time"
)

var (
	GeneratorMap      = make(map[uint][]*corev1.Pod)
	GeneratorMapMutex sync.RWMutex
)

// StartGenerator 启动动态附件生成器, 等待附加命令, 生成附件, contestChallenge 需要预加载 Challenge
func StartGenerator(contestChallenge model.ContestChallenge) (*corev1.Pod, bool, string) {
	var (
		pod           *corev1.Pod
		ok            bool
		msg           string
		err           error
		generatorName = fmt.Sprintf("gen-%s", utils.RandStr(20))
		containerName = fmt.Sprintf("ctn-%s", utils.RandStr(20))
		volumeName    = fmt.Sprintf("vol-%s", utils.RandStr(20))
		lables        = map[string]string{GeneratorPodTag: generatorName, "contest_challenge_id": fmt.Sprintf("%d", contestChallenge.ID)}
	)
	if contestChallenge.Challenge.GeneratorImage == "" {
		return nil, false, i18n.InvalidDockerImage
	}
	log.Logger.Infof("Starting Generator for Challenge %d-%s", contestChallenge.ChallengeID, contestChallenge.Name)
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()
	pwd := utils.UUID()
	pod, ok, msg = CreatePod(ctx, CreatePodOptions{
		Name:   generatorName,
		Labels: lables,
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
				VolumeMounts: []corev1.VolumeMount{
					{
						Name:      volumeName,
						MountPath: "/root/mnt",
						SubPath: strings.TrimPrefix(
							strings.TrimPrefix(contestChallenge.Challenge.BasicDir(), config.Env.Path), "/",
						),
					},
				},
				WorkingDir: "/root",
				Command:    []string{"sleep", "infinity"},
			},
		},
		Volumes: []corev1.Volume{
			{
				Name: volumeName,
				VolumeSource: corev1.VolumeSource{
					PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
						ClaimName: NFSVolumeName,
					},
				},
			},
		},
	})
	if !ok {
		return nil, false, msg
	}
	var commands []string
	if _, err = os.Stat(contestChallenge.Challenge.GeneratorPath()); err == nil {
		commands = append(commands, fmt.Sprintf("unzip /root/mnt/generator.zip -d /root"))
	} else {
		log.Logger.Info("Generator file not found, make sure the generator docker can work correctly")
	}
	for _, command := range commands {
		log.Logger.Debugf("Executing command: %s", command)
		if _, _, err = Exec(generatorName, containerName, command, nil); err != nil {
			log.Logger.Warningf("Failed to execute command %s: %v", command, err)
			return nil, false, i18n.ExecCommandError
		}
	}
	GeneratorMapMutex.Lock()
	defer GeneratorMapMutex.Unlock()
	GeneratorMap[contestChallenge.ID] = append(GeneratorMap[contestChallenge.ID], pod)
	return pod, true, i18n.Success
}

func GetGenerator(contestChallenge model.ContestChallenge) (*corev1.Pod, bool, string) {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	GeneratorMapMutex.Lock()
	defer GeneratorMapMutex.Unlock()
	generators, ok := GeneratorMap[contestChallenge.ID]
	if !ok {
		return StartGenerator(contestChallenge)
	}
	if len(generators) > 0 {
		index := rand.Intn(len(generators))
		return generators[index], true, i18n.Success
	}
	return nil, false, i18n.UnknownError
}

// StopGenerator 停止动态附件生成器, contestChallenge 需要预加载 Challenge
func StopGenerator(contestChallenge model.ContestChallenge, generator *corev1.Pod) (bool, string) {
	log.Logger.Infof("Stopping generator for challenge %d-%s", contestChallenge.ChallengeID, contestChallenge.Name)

	GeneratorMapMutex.Lock()
	defer GeneratorMapMutex.Unlock()
	_, ok := GeneratorMap[contestChallenge.ID]
	if ok {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
		defer cancel()
		if ok, msg := DeletePod(ctx, generator.Name); !ok {
			return false, msg
		}
		labels := map[string]string{GeneratorPodTag: generator.Name, "contest_challenge_id": fmt.Sprintf("%d", contestChallenge.ID)}
		if ok, msg := DeleteServiceList(ctx, labels); !ok {
			return false, msg
		}
		GeneratorMap[contestChallenge.ID] = slices.DeleteFunc(GeneratorMap[contestChallenge.ID], func(pod *corev1.Pod) bool {
			return pod.Name == generator.Name
		})
	}
	return true, i18n.Success
}

// GenerateAttachment 附加容器命令, 生成附件, model.Usage 需要预加载
func GenerateAttachment(contestChallenge model.ContestChallenge, team model.Team, teamFlagL []model.TeamFlag) (bool, string) {
	var err error
	log.Logger.Debugf("Generating attachment for team %d challenge %d", team.ID, contestChallenge.ChallengeID)
	generator, ok, msg := GetGenerator(contestChallenge)
	// 附加失败则直接返回, 并尝试关闭生成器
	if !ok || generator.Status.Phase != corev1.PodRunning {
		go StopGenerator(contestChallenge, generator)
		return false, msg
	}
	var flags string
	for _, teamFlag := range teamFlagL {
		flags += fmt.Sprintf("%s,", base64.StdEncoding.EncodeToString([]byte(teamFlag.Value)))
	}
	flags = base64.StdEncoding.EncodeToString([]byte(strings.TrimSuffix(flags, ",")))
	flags = strings.TrimSuffix(flags, ",")
	command := fmt.Sprintf("./run.sh %d %s", team.ID, base64.StdEncoding.EncodeToString([]byte(flags)))
	log.Logger.Debugf("Executing command in %s: %s", generator.Name, command)
	if _, _, err = Exec(generator.Name, generator.Spec.Containers[0].Name, command, nil); err != nil {
		log.Logger.Warningf("Failed to execute command %s: %v", command, err)
		return false, i18n.ExecCommandError
	}
	return true, i18n.Success
}
