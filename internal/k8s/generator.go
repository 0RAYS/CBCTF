package k8s

import (
	"CBCTF/internal/config"
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"CBCTF/internal/utils"
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"math/big"
	"os"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"

	corev1 "k8s.io/api/core/v1"
)

var (
	GeneratorMap      = make(map[uint][]*Generator)
	GeneratorMapMutex sync.RWMutex
)

type Generator struct {
	Start time.Time
	Pod   *corev1.Pod
}

// StartGenerator 启动动态附件生成器, 等待附加命令, 生成附件
func StartGenerator(ctx context.Context, challenge model.Challenge) (*corev1.Pod, bool, string) {
	var (
		pod           *corev1.Pod
		ok            bool
		msg           string
		err           error
		generatorName = fmt.Sprintf("gen-%s", utils.RandStr(20))
		containerName = fmt.Sprintf("ctn-%s", utils.RandStr(20))
		volumeName    = fmt.Sprintf("vol-%s", utils.RandStr(20))
		labels        = map[string]string{GeneratorPodTag: generatorName, "contest_challenge_id": strconv.Itoa(int(challenge.ID))}
	)
	if challenge.GeneratorImage == "" {
		return nil, false, i18n.InvalidDockerImage
	}
	log.Logger.Infof("Starting Generator for Challenge %d-%s", challenge.ID, challenge.Name)
	pwd := utils.UUID()
	pod, ok, msg = CreatePod(ctx, CreatePodOptions{
		Name:   generatorName,
		Labels: labels,
		Containers: []corev1.Container{
			{
				Name:  containerName,
				Image: challenge.GeneratorImage,
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
							strings.TrimPrefix(challenge.BasicDir(), config.Env.Path), "/",
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
						ClaimName: nfsVolumeName,
					},
				},
			},
		},
	})
	if !ok {
		return nil, false, msg
	}
	var commands []string
	if _, err = os.Stat(challenge.GeneratorPath()); err == nil {
		commands = append(commands, fmt.Sprintf("unzip /root/mnt/generator.zip -d /root"))
	} else {
		log.Logger.Info("Generator file not found, make sure the generator docker can work correctly")
	}
	for _, command := range commands {
		log.Logger.Debugf("Executing command: %s", command)
		if _, _, err = Exec(generatorName, containerName, command, nil); err != nil {
			log.Logger.Warningf("Failed to execute command %s: %s", command, err)
			return nil, false, i18n.ExecCommandError
		}
	}
	GeneratorMapMutex.Lock()
	GeneratorMap[challenge.ID] = append(GeneratorMap[challenge.ID], &Generator{Start: time.Now(), Pod: pod})
	GeneratorMapMutex.Unlock()
	return pod, true, i18n.Success
}

func GetGenerator(ctx context.Context, challenge model.Challenge) (*corev1.Pod, bool, string) {
	GeneratorMapMutex.RLock()
	generators, ok := GeneratorMap[challenge.ID]
	GeneratorMapMutex.RUnlock()
	if !ok {
		return StartGenerator(ctx, challenge)
	}
	if len(generators) > 0 {
		index, _ := rand.Int(rand.Reader, big.NewInt(int64(len(generators))))
		return generators[index.Int64()].Pod, true, i18n.Success
	}
	return nil, false, i18n.UnknownError
}

// StopGenerator 停止动态附件生成器
func StopGenerator(ctx context.Context, challenge model.Challenge, generator *corev1.Pod) (bool, string) {
	log.Logger.Infof("Stopping generator for Challenge %d-%s", challenge.ID, challenge.Name)
	GeneratorMapMutex.RLock()
	_, ok := GeneratorMap[challenge.ID]
	GeneratorMapMutex.RUnlock()
	if ok {
		if ok, msg := DeletePod(ctx, generator.Name); !ok {
			return false, msg
		}
		labels := map[string]string{GeneratorPodTag: generator.Name, "contest_challenge_id": strconv.Itoa(int(challenge.ID))}
		if ok, msg := DeleteServiceList(ctx, labels); !ok {
			return false, msg
		}
		GeneratorMapMutex.Lock()
		GeneratorMap[challenge.ID] = slices.DeleteFunc(GeneratorMap[challenge.ID], func(gen *Generator) bool {
			return gen.Pod.Name == generator.Name
		})
		GeneratorMapMutex.Unlock()
	}
	return true, i18n.Success
}

// GenAttachment 附加容器命令, 生成附件
func GenAttachment(ctx context.Context, challenge model.Challenge, team model.Team, teamFlagL []model.TeamFlag) (bool, string) {
	var err error
	log.Logger.Debugf("Generating attachment for Team %d Challenge %d", team.ID, challenge.ID)
	generator, ok, _ := GetGenerator(ctx, challenge)
	// 附加失败则直接返回, 并尝试关闭生成器
	if !ok || generator.Status.Phase != corev1.PodRunning {
		return StopGenerator(ctx, challenge, generator)
	}
	var flags string
	for _, teamFlag := range teamFlagL {
		flags += fmt.Sprintf("%s,", base64.StdEncoding.EncodeToString([]byte(teamFlag.Value)))
	}
	flags = base64.StdEncoding.EncodeToString([]byte(strings.TrimSuffix(flags, ",")))
	flags = strings.TrimSuffix(flags, ",")
	_ = os.Remove(challenge.AttachmentPath(team.ID))
	command := fmt.Sprintf("./run.sh %d %s", team.ID, base64.StdEncoding.EncodeToString([]byte(flags)))
	log.Logger.Debugf("Executing command in %s: %s", generator.Name, command)
	if _, _, err = Exec(generator.Name, generator.Spec.Containers[0].Name, command, nil); err != nil {
		log.Logger.Warningf("Failed to execute command %s: %s", command, err)
		return false, i18n.ExecCommandError
	}
	for {
		if _, err = os.Stat(challenge.AttachmentPath(team.ID)); err == nil {
			break
		}
		time.Sleep(time.Second)
	}
	return true, i18n.Success
}

func GenTestAttachment(ctx context.Context, challenge model.Challenge, challengeFlags []model.ChallengeFlag) (bool, string) {
	var err error
	log.Logger.Debugf("Generating test attachment for Challenge %d", challenge.ID)
	generator, ok, _ := GetGenerator(ctx, challenge)
	// 附加失败则直接返回, 并尝试关闭生成器
	if !ok || generator.Status.Phase != corev1.PodRunning {
		return StopGenerator(ctx, challenge, generator)
	}
	var flags string
	for _, flag := range challengeFlags {
		flags += fmt.Sprintf("%s,", base64.StdEncoding.EncodeToString([]byte(flag.Value)))
	}
	flags = base64.StdEncoding.EncodeToString([]byte(strings.TrimSuffix(flags, ",")))
	flags = strings.TrimSuffix(flags, ",")
	_ = os.Remove(challenge.AttachmentPath(0))
	command := fmt.Sprintf("./run.sh %d %s", 0, base64.StdEncoding.EncodeToString([]byte(flags)))
	log.Logger.Debugf("Executing command in %s: %s", generator.Name, command)
	if _, _, err = Exec(generator.Name, generator.Spec.Containers[0].Name, command, nil); err != nil {
		log.Logger.Warningf("Failed to execute command %s: %s", command, err)
		return false, i18n.ExecCommandError
	}
	for {
		if _, err = os.Stat(challenge.AttachmentPath(0)); err == nil {
			break
		}
		time.Sleep(time.Second)
	}
	return true, i18n.Success
}
