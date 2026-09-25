package k8s

import (
	"CBCTF/internal/config"
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"CBCTF/internal/worker"
	"archive/zip"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"maps"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func GeneratorLabels(generator model.Generator, tags ...map[string]string) map[string]string {
	labels := map[string]string{
		"challenge_id": strconv.Itoa(int(generator.ChallengeID)),
		"generator_id": strconv.Itoa(int(generator.ID)),
	}
	if len(tags) > 0 {
		maps.Copy(labels, tags[0])
	}
	return labels
}

func StartGenerator(ctx context.Context, challenge model.Challenge, generator model.Generator) (*corev1.Pod, model.RetVal) {
	var (
		pod    *corev1.Pod
		ret    model.RetVal
		labels = GeneratorLabels(generator, map[string]string{RoleLabel: GeneratorPodTag})
	)
	if challenge.GeneratorImage == "" {
		return nil, model.RetVal{Msg: i18n.Model.Challenge.EmptyImage}
	}
	log.Logger.Debugf("Creating generator pod: generator_id=%d name=%s challenge_id=%d image=%s", generator.ID, generator.Name, challenge.ID, challenge.GeneratorImage)
	pod, ret = CreatePod(ctx, CreatePodOptions{
		PriorityClassName: config.Env.K8S.PriorityClassName,
		Name:              generator.Name,
		Labels:            labels,
		InitContainers: []corev1.Container{{
			Name: "worker-install", Image: config.Env.K8S.WorkerImage, ImagePullPolicy: corev1.PullIfNotPresent,
			Command:      []string{"/app/worker", "--install", "/worker/worker"},
			VolumeMounts: []corev1.VolumeMount{{Name: "worker", MountPath: "/worker"}}, Resources: sidecarResources(),
		}},
		Containers: []corev1.Container{
			{
				Name:            "generator",
				Resources:       corev1.ResourceRequirements{Requests: workloadRequests(nil)},
				Image:           challenge.GeneratorImage,
				ImagePullPolicy: corev1.PullIfNotPresent,
				Env:             []corev1.EnvVar{{Name: "CBCTF_WORKER_TOKEN", Value: generator.WorkerToken}},
				ReadinessProbe:  &corev1.Probe{ProbeHandler: corev1.ProbeHandler{HTTPGet: &corev1.HTTPGetAction{Path: "/ready", Port: intstr.FromInt(worker.Port)}}, PeriodSeconds: 1, TimeoutSeconds: 1},
				VolumeMounts: []corev1.VolumeMount{
					{Name: "worker", MountPath: "/worker", ReadOnly: true},
					{Name: "output", MountPath: "/root/mnt/attachments"},
					{
						Name:      nfsVolumeName,
						MountPath: "/root/mnt",
						SubPath: strings.TrimPrefix(
							strings.TrimPrefix(challenge.BasicDir(), config.Env.Path), "/",
						),
					},
				},
				WorkingDir: "/root",
				Command:    []string{"/worker/worker"},
			},
		},
		Volumes: []corev1.Volume{
			{Name: "worker", EmptyDir: &corev1.EmptyDirVolumeSource{}},
			{Name: "output", EmptyDir: &corev1.EmptyDirVolumeSource{}},
			{
				Name: nfsVolumeName,
				PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
					ClaimName: nfsVolumeName,
				},
			},
		},
	})
	if !ret.OK {
		return nil, ret
	}
	log.Logger.Debugf("Created generator pod: generator_id=%d name=%s challenge_id=%d", generator.ID, generator.Name, challenge.ID)
	return pod, model.SuccessRetVal()
}

func StopGenerator(ctx context.Context, generator model.Generator) model.RetVal {
	log.Logger.Debugf("Deleting generator k8s resources: generator_id=%d name=%s challenge_id=%d", generator.ID, generator.Name, generator.ChallengeID)
	pod, ret := GetPod(ctx, generator.Name)
	if !ret.OK && ret.Msg != i18n.K8S.NotFound {
		return ret
	}
	if ret.OK {
		if !GeneratorOwnsPod(generator, pod) {
			return model.RetVal{Msg: i18n.K8S.DeleteError, Attr: map[string]any{"Model": "Pod", "Error": "refusing to delete a Pod not owned by this generator"}}
		}
		if ret := DeletePodAndWait(ctx, generator.Name, pod.UID); !ret.OK {
			return ret
		}
	}
	labels := GeneratorLabels(generator, map[string]string{RoleLabel: GeneratorPodTag})
	if ret := DeleteServiceCollection(ctx, labels); !ret.OK {
		return ret
	}
	log.Logger.Debugf("Deleted generator k8s resources: generator_id=%d name=%s challenge_id=%d", generator.ID, generator.Name, generator.ChallengeID)
	return model.SuccessRetVal()
}

// GenAttachment streams a completed artifact directly from the worker into an
// immutable cache path. The worker's output directory is local EmptyDir, so
// completion never depends on an NFS client's negative/attribute cache.
func GenAttachment(ctx context.Context, challenge model.Challenge, generator model.Generator, teamID uint, flags []string) model.RetVal {
	log.Logger.Debugf("Running attachment generator: team_id=%d challenge_id=%d generator_id=%d", teamID, challenge.ID, generator.ID)
	if err := generateAttachment(ctx, challenge, generator, teamID, flags); err != nil {
		if errors.Is(err, errGeneratorMissing) {
			return model.RetVal{Msg: i18n.K8S.NotFound, Attr: map[string]any{"Model": "Generator"}}
		}
		return model.RetVal{Msg: i18n.Common.UnknownError, Attr: map[string]any{"Error": err.Error()}}
	}
	return model.SuccessRetVal()
}

var workerHTTP = &http.Client{Timeout: 90 * time.Second, Transport: &http.Transport{Proxy: nil, MaxIdleConns: 100, MaxIdleConnsPerHost: 2, IdleConnTimeout: time.Minute}}
var errGeneratorMissing = errors.New("generator pod is missing")

func generateAttachment(ctx context.Context, challenge model.Challenge, generator model.Generator, teamID uint, flags []string) error {
	pod, err := cachedPod(ctx, generator.Name)
	if err != nil {
		return err
	}
	if pod == nil {
		return errGeneratorMissing
	}
	if !GeneratorOwnsPod(generator, pod) || !podReady(pod) {
		return fmt.Errorf("generator is not ready")
	}
	if generator.WorkerToken == "" {
		return fmt.Errorf("generator has no worker token")
	}
	revision := worker.SourceRevision(challenge.GeneratorPath())
	destination := challenge.AttachmentCachePathForRevision(teamID, flags, revision)
	body, err := json.Marshal(worker.Request{TeamID: teamID, Flags: flags, SourceRevision: revision})
	if err != nil {
		return err
	}
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, "http://"+net.JoinHostPort(pod.Status.PodIP, strconv.Itoa(worker.Port))+"/generate", bytes.NewReader(body))
	if err != nil {
		return err
	}
	request.Header.Set("Authorization", "Bearer "+generator.WorkerToken)
	request.Header.Set("Content-Type", "application/json")
	response, err := workerHTTP.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("generator returned HTTP %d", response.StatusCode)
	}
	if err := os.MkdirAll(filepath.Dir(destination), 0700); err != nil {
		return err
	}
	file, err := os.CreateTemp(filepath.Dir(destination), ".attachment-*")
	if err != nil {
		return err
	}
	defer os.Remove(file.Name())
	defer file.Close()
	if _, err = io.Copy(file, response.Body); err != nil {
		return err
	}
	if err = file.Sync(); err != nil {
		return err
	}
	if err = file.Close(); err != nil {
		return err
	}
	archive, err := zip.OpenReader(file.Name())
	if err != nil {
		return fmt.Errorf("invalid generated ZIP: %w", err)
	}
	if err = archive.Close(); err != nil {
		return err
	}
	return os.Rename(file.Name(), destination)
}
