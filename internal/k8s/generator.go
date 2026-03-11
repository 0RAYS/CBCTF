package k8s

import (
	"CBCTF/internal/config"
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	corev1 "k8s.io/api/core/v1"
)

func StartGenerator(ctx context.Context, challenge model.Challenge, generator model.Generator) (*corev1.Pod, model.RetVal) {
	var (
		pod    *corev1.Pod
		ret    model.RetVal
		err    error
		labels = map[string]string{
			GeneratorPodTag: GeneratorPodTag,
			"challenge_id":  strconv.Itoa(int(challenge.ID)),
		}
	)
	if challenge.GeneratorImage == "" {
		return nil, model.RetVal{Msg: i18n.Model.Challenge.EmptyImage}
	}
	log.Logger.Infof("Starting Generator %s for Challenge %d-%s", generator.Name, challenge.ID, challenge.Name)
	pod, ret = CreatePod(ctx, CreatePodOptions{
		Name:   generator.Name,
		Labels: labels,
		Containers: []corev1.Container{
			{
				Name:  "generator",
				Image: challenge.GeneratorImage,
				VolumeMounts: []corev1.VolumeMount{
					{
						Name:      nfsVolumeName,
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
				Name: nfsVolumeName,
				VolumeSource: corev1.VolumeSource{
					PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
						ClaimName: nfsVolumeName,
					},
				},
			},
		},
	})
	if !ret.OK {
		return nil, ret
	}
	var commands []string
	if _, err = os.Stat(challenge.GeneratorPath()); err == nil {
		commands = append(commands, fmt.Sprintf("unzip /root/mnt/generator.zip -d /root"))
	} else {
		log.Logger.Info("Generator file not found, make sure the generator docker can work correctly")
	}
	for _, command := range commands {
		log.Logger.Debugf("Executing command: %s", command)
		if _, _, err = Exec(ctx, generator.Name, pod.Spec.Containers[0].Name, command, nil); err != nil {
			log.Logger.Warningf("Failed to execute command %s: %s", command, err)
			return nil, model.RetVal{Msg: i18n.Common.UnknownError, Attr: map[string]any{"Error": err.Error()}}
		}
	}
	return pod, model.SuccessRetVal()
}

func StopGenerator(ctx context.Context, generator model.Generator) model.RetVal {
	log.Logger.Infof("Stopping generator %s for Challenge %d", generator.Name, generator.ChallengeID)
	if ret := DeletePod(ctx, generator.Name); !ret.OK {
		return ret
	}
	labels := map[string]string{
		GeneratorPodTag: generator.Name,
		"challenge_id":  strconv.Itoa(int(generator.ChallengeID)),
	}
	if ret := DeleteServiceList(ctx, labels); !ret.OK {
		return ret
	}
	return model.SuccessRetVal()
}

// GenAttachment 附加容器命令, 生成附件
func GenAttachment(ctx context.Context, challenge model.Challenge, generator model.Generator, teamID uint, flags []string) model.RetVal {
	var err error
	log.Logger.Debugf("Generating attachment for Team %d Challenge %d", teamID, challenge.ID)
	pod, ret := GetPod(ctx, generator.Name)
	if !ret.OK {
		return ret
	}
	if pod.Status.Phase != corev1.PodRunning {
		StopGenerator(ctx, generator)
		return model.RetVal{Msg: i18n.Model.NotFound, Attr: map[string]any{"Model": generator.ModelName()}}
	}
	var flag string
	for _, value := range flags {
		flag += fmt.Sprintf("%s,", base64.StdEncoding.EncodeToString([]byte(value)))
	}
	flag = base64.StdEncoding.EncodeToString([]byte(strings.TrimSuffix(flag, ",")))
	filepath := challenge.AttachmentPath(teamID)
	_ = os.Remove(filepath)
	command := fmt.Sprintf("./run.sh %d %s", teamID, flag)
	log.Logger.Debugf("Executing command in %s: %s", generator.Name, command)
	if _, _, err = Exec(ctx, generator.Name, pod.Spec.Containers[0].Name, command, nil); err != nil {
		log.Logger.Warningf("Failed to execute command %s: %s", command, err)
		return model.RetVal{Msg: i18n.Common.UnknownError, Attr: map[string]any{"Error": err.Error()}}
	}
	for {
		select {
		case <-ctx.Done():
			return model.RetVal{Msg: i18n.Common.UnknownError, Attr: map[string]any{"Error": "timeout waiting for attachment"}}
		default:
		}
		// NFS 延迟写入, 主动触发读取
		_, _ = os.ReadDir(path.Dir(filepath))
		if _, err = os.Stat(filepath); err == nil {
			break
		}
		time.Sleep(time.Second)
	}
	return model.SuccessRetVal()

}
