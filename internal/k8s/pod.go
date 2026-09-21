package k8s

import (
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"context"
	"fmt"
	"io"
	"maps"
	"strings"
	"time"

	corev1 "k8s.io/api/core/v1"
	apierror "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
)

type Network struct {
	Interface    string
	IPv4         string
	MAC          string
	Subnet       string
	NetAttachDef string
}

type CreatePodOptions struct {
	Name        string
	Labels      map[string]string
	Annotations map[string]string
	Networks    []Network
	Containers  []corev1.Container
	Volumes     []corev1.Volume
}

func CreatePod(ctx context.Context, options CreatePodOptions) (*corev1.Pod, model.RetVal) {
	var (
		pod *corev1.Pod
		err error
	)
	pod = &corev1.Pod{
		Name:      options.Name,
		Namespace: globalNamespace,
		Labels:    options.Labels,
		Annotations: func() map[string]string {
			annotations := make(map[string]string)
			maps.Copy(annotations, options.Annotations)
			for _, network := range options.Networks {
				annotations["k8s.v1.cni.cncf.io/networks"] += fmt.Sprintf(",%s/%s", globalNamespace, network.NetAttachDef)
				annotations["k8s.v1.cni.cncf.io/networks"] = strings.Trim(annotations["k8s.v1.cni.cncf.io/networks"], ",")
				annotations[fmt.Sprintf("%s.%s.ovn.kubernetes.io/logical_switch", network.NetAttachDef, globalNamespace)] = network.Subnet
				annotations[fmt.Sprintf("%s.%s.ovn.kubernetes.io/ip_address", network.NetAttachDef, globalNamespace)] = network.IPv4
				if network.MAC != "" {
					annotations[fmt.Sprintf("%s.%s.ovn.kubernetes.io/mac_address", network.NetAttachDef, globalNamespace)] = network.MAC
				}
			}
			if len(annotations) == 0 {
				return nil
			}
			return annotations
		}(),
		Spec: corev1.PodSpec{
			EnableServiceLinks:            new(false),
			AutomountServiceAccountToken:  new(false),
			Containers:                    options.Containers,
			Volumes:                       options.Volumes,
			TerminationGracePeriodSeconds: new(int64(3)),
			RestartPolicy:                 corev1.RestartPolicyNever,
		},
	}
	pod, err = kubeClient.CoreV1().Pods(globalNamespace).Create(ctx, pod, metav1.CreateOptions{})
	if err != nil {
		log.Logger.Warningf("Failed to create Pod: %s", err)
		return nil, model.RetVal{Msg: i18n.K8S.CreateError, Attr: map[string]any{"Model": "Pod", "Error": err.Error()}}
	}
	ready, err := waitPodReady(ctx, pod)
	if err != nil {
		return nil, model.RetVal{Msg: i18n.K8S.PodRunError, Attr: map[string]any{
			"Pod": pod.Name, "Phase": pod.Status.Phase, "Reason": err.Error(), "Error": err.Error(),
		}}
	}
	return ready, model.SuccessRetVal()
}

// GetPod 依据 name 获取 Pod
func GetPod(ctx context.Context, name string) (*corev1.Pod, model.RetVal) {
	pod, err := kubeClient.CoreV1().Pods(globalNamespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		if apierror.IsNotFound(err) {
			return nil, model.RetVal{Msg: i18n.K8S.NotFound, Attr: map[string]any{"Model": "Pod"}}
		}
		log.Logger.Warningf("Failed to get Pod %s: %s", name, err)
		return nil, model.RetVal{Msg: i18n.K8S.GetError, Attr: map[string]any{"Model": "Pod", "Error": err.Error()}}
	}
	return pod, model.SuccessRetVal()
}

func ListPods(ctx context.Context, labels ...map[string]string) (*corev1.PodList, model.RetVal) {
	var options metav1.ListOptions
	if len(labels) > 0 {
		var selector strings.Builder
		for k, v := range labels[0] {
			selector.WriteString(fmt.Sprintf("%s=%s,", k, v))
		}
		options = metav1.ListOptions{
			LabelSelector: strings.TrimSuffix(selector.String(), ","),
		}
	}
	podList, err := kubeClient.CoreV1().Pods(globalNamespace).List(ctx, options)
	if err != nil {
		log.Logger.Warningf("Failed to list Pods: %s", err)
		return nil, model.RetVal{Msg: i18n.K8S.GetError, Attr: map[string]any{"Model": "Pod", "Error": err.Error()}}
	}
	return podList, model.SuccessRetVal()
}

const MaxPodLogLines int64 = 10000
const MaxPodLogBytes int64 = 1 << 20

func podLogOptions(containerName string, lines int64) *corev1.PodLogOptions {
	if lines <= 0 {
		lines = 1000
	}
	lines = min(lines, MaxPodLogLines)
	return &corev1.PodLogOptions{Container: containerName, TailLines: &lines, LimitBytes: new(MaxPodLogBytes)}
}

func GetPodLogs(ctx context.Context, podName, containerName string, lines int64) (string, model.RetVal) {
	options := podLogOptions(containerName, lines)
	podLogs, err := kubeClient.CoreV1().Pods(globalNamespace).GetLogs(podName, options).Stream(ctx)
	if err != nil {
		log.Logger.Warningf("Failed to get Pod Logs: %s", err)
		return "", model.RetVal{Msg: i18n.K8S.GetError, Attr: map[string]any{"Model": "PodLog", "Error": err.Error()}}
	}
	defer func(podLogs io.ReadCloser) {
		_ = podLogs.Close()
	}(podLogs)
	buf, err := io.ReadAll(io.LimitReader(podLogs, MaxPodLogBytes))
	if err != nil {
		log.Logger.Warningf("Failed to read Pod Logs: %s", err)
		return "", model.RetVal{Msg: i18n.K8S.GetError, Attr: map[string]any{"Model": "PodLog", "Error": err.Error()}}
	}
	return string(buf), model.SuccessRetVal()
}

// DeletePod 依据 name 删除 Pod
func DeletePod(ctx context.Context, name string) model.RetVal {
	err := kubeClient.CoreV1().Pods(globalNamespace).Delete(ctx, name, metav1.DeleteOptions{})
	if err != nil && !apierror.IsNotFound(err) {
		log.Logger.Warningf("Failed to delete Pod: %s", err)
		return model.RetVal{Msg: i18n.K8S.DeleteError, Attr: map[string]any{"Model": "Pod", "Error": err.Error()}}
	}
	return model.SuccessRetVal()
}

func DeletePodCollection(ctx context.Context, labels ...map[string]string) model.RetVal {
	var options metav1.ListOptions
	if len(labels) > 0 {
		var selector strings.Builder
		for k, v := range labels[0] {
			selector.WriteString(fmt.Sprintf("%s=%s,", k, v))
		}
		options = metav1.ListOptions{
			LabelSelector: strings.TrimSuffix(selector.String(), ","),
		}
	}
	err := kubeClient.CoreV1().Pods(globalNamespace).DeleteCollection(ctx, metav1.DeleteOptions{}, options)
	if err != nil && !apierror.IsNotFound(err) {
		log.Logger.Warningf("Failed to delete Pod: %s", err)
		return model.RetVal{Msg: i18n.K8S.DeleteError, Attr: map[string]any{"Model": "Pod", "Error": err.Error()}}
	}
	return model.SuccessRetVal()
}

// Delete acceptance is not deletion completion (PVC mounts and exec may still be active).
func DeletePodAndWait(ctx context.Context, name string) model.RetVal {
	if ret := DeletePod(ctx, name); !ret.OK {
		return ret
	}
	err := wait.PollUntilContextCancel(ctx, 500*time.Millisecond, true, func(ctx context.Context) (bool, error) {
		_, err := kubeClient.CoreV1().Pods(globalNamespace).Get(ctx, name, metav1.GetOptions{})
		if apierror.IsNotFound(err) {
			return true, nil
		}
		return false, err
	})
	if err != nil {
		return model.RetVal{Msg: i18n.K8S.DeleteError, Attr: map[string]any{"Model": "Pod", "Error": err.Error()}}
	}
	return model.SuccessRetVal()
}
