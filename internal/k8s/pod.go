package k8s

import (
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	corev1 "k8s.io/api/core/v1"
	apierror "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Network struct {
	Interface    string
	IPv4         string
	MAC          string
	Gateway      string
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
		ret model.RetVal
		err error
	)
	if _, ret = GetPod(ctx, options.Name); !ret.OK {
		DeletePod(ctx, options.Name)
	}
	pod = &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      options.Name,
			Namespace: globalNamespace,
			Labels:    options.Labels,
			Annotations: func() map[string]string {
				annotations := make(map[string]string)
				for key, value := range options.Annotations {
					annotations[key] = value
				}
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
		},
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
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()
	for {
		pod, ret = GetPod(ctx, options.Name)
		if !ret.OK {
			return nil, ret
		}
		if pod.Status.Phase == corev1.PodRunning {
			break
		}
		if pod.Status.Phase != corev1.PodPending {
			log.Logger.Warningf("Failed to run Pod %s: phase=%s, reason=%s", pod.Name, pod.Status.Phase, pod.Status.Reason)
			return nil, model.RetVal{Msg: i18n.K8S.PodRunError, Attr: map[string]any{"Pod": pod.Name, "Phase": pod.Status.Phase, "Reason": pod.Status.Reason}}
		}
		<-ticker.C
	}
	return pod, model.SuccessRetVal()
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
		var selector string
		for k, v := range labels[0] {
			selector += fmt.Sprintf("%s=%s,", k, v)
		}
		options = metav1.ListOptions{
			LabelSelector: strings.TrimSuffix(selector, ","),
		}
	}
	podList, err := kubeClient.CoreV1().Pods(globalNamespace).List(ctx, options)
	if err != nil {
		log.Logger.Warningf("Failed to list Pods: %s", err)
		return nil, model.RetVal{Msg: i18n.K8S.GetError, Attr: map[string]any{"Model": "Pod", "Error": err.Error()}}
	}
	return podList, model.SuccessRetVal()
}

func GetPodLogs(ctx context.Context, podName, containerName string, lines int64) (string, model.RetVal) {
	options := &corev1.PodLogOptions{
		Container: containerName,
		Follow:    false,
		TailLines: &lines,
	}
	podLogs, err := kubeClient.CoreV1().Pods(globalNamespace).GetLogs(podName, options).Stream(ctx)
	if err != nil {
		log.Logger.Warningf("Failed to get Pod Logs: %s", err)
		return "", model.RetVal{Msg: i18n.K8S.GetError, Attr: map[string]any{"Model": "PodLog", "Error": err.Error()}}
	}
	defer func(podLogs io.ReadCloser) {
		_ = podLogs.Close()
	}(podLogs)
	buf, err := io.ReadAll(podLogs)
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
		var selector string
		for k, v := range labels[0] {
			selector += fmt.Sprintf("%s=%s,", k, v)
		}
		options = metav1.ListOptions{
			LabelSelector: strings.TrimSuffix(selector, ","),
		}
	}
	err := kubeClient.CoreV1().Pods(globalNamespace).DeleteCollection(ctx, metav1.DeleteOptions{}, options)
	if err != nil && !apierror.IsNotFound(err) {
		log.Logger.Warningf("Failed to delete Pod: %s", err)
		return model.RetVal{Msg: i18n.K8S.DeleteError, Attr: map[string]any{"Model": "Pod", "Error": err.Error()}}
	}
	return model.SuccessRetVal()
}
