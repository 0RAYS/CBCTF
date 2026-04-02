package k8s

import (
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"context"
	"fmt"
	"strings"
	"time"

	corev1 "k8s.io/api/core/v1"
	apierror "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type CreatePodOptions struct {
	Name            string
	Labels          map[string]string
	Annotations     map[string]string
	Containers      []corev1.Container
	Volumes         []corev1.Volume
	PodAffinity     map[string]string
	PodAntiAffinity map[string]string
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
			Name:        options.Name,
			Namespace:   globalNamespace,
			Labels:      options.Labels,
			Annotations: options.Annotations,
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
	if len(options.PodAffinity) > 0 || len(options.PodAntiAffinity) > 0 {
		pod.Spec.Affinity = &corev1.Affinity{}
		if len(options.PodAffinity) > 0 {
			pod.Spec.Affinity.PodAffinity = &corev1.PodAffinity{
				RequiredDuringSchedulingIgnoredDuringExecution: []corev1.PodAffinityTerm{
					{
						LabelSelector: &metav1.LabelSelector{
							MatchExpressions: func() []metav1.LabelSelectorRequirement {
								tmp := make([]metav1.LabelSelectorRequirement, 0)
								for key, value := range options.PodAffinity {
									tmp = append(tmp, metav1.LabelSelectorRequirement{
										Key:      key,
										Operator: metav1.LabelSelectorOpIn,
										Values:   []string{value},
									})
								}
								return tmp
							}(),
						},
						Namespaces:  []string{globalNamespace, "kube-system"},
						TopologyKey: "kubernetes.io/hostname",
					},
				},
			}
		}
		if len(options.PodAntiAffinity) > 0 {
			pod.Spec.Affinity.PodAntiAffinity = &corev1.PodAntiAffinity{
				RequiredDuringSchedulingIgnoredDuringExecution: []corev1.PodAffinityTerm{
					{
						LabelSelector: &metav1.LabelSelector{
							MatchExpressions: func() []metav1.LabelSelectorRequirement {
								tmp := make([]metav1.LabelSelectorRequirement, 0)
								for key, value := range options.PodAntiAffinity {
									tmp = append(tmp, metav1.LabelSelectorRequirement{
										Key:      key,
										Operator: metav1.LabelSelectorOpIn,
										Values:   []string{value},
									})
								}
								return tmp
							}(),
						},
						Namespaces:  []string{globalNamespace, "kube-system"},
						TopologyKey: "kubernetes.io/hostname",
					},
				},
			}
		}
	}
	pod, err = kubeClient.CoreV1().Pods(globalNamespace).Create(ctx, pod, metav1.CreateOptions{})
	if err != nil {
		log.Logger.Warningf("Failed to create Pod: %s", err)
		return nil, model.RetVal{Msg: i18n.K8S.CreateError, Attr: map[string]any{"Model": "Pod", "Error": err.Error()}}
	}
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
			return nil, model.RetVal{Msg: i18n.Common.UnknownError, Attr: map[string]any{"Error": "Failed to run Pod"}}
		}
		time.Sleep(500 * time.Millisecond)
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

func GetPodList(ctx context.Context, labels ...map[string]string) (*corev1.PodList, model.RetVal) {
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

// DeletePod 依据 name 删除 Pod
func DeletePod(ctx context.Context, name string) model.RetVal {
	err := kubeClient.CoreV1().Pods(globalNamespace).Delete(ctx, name, metav1.DeleteOptions{})
	if err != nil && !apierror.IsNotFound(err) {
		log.Logger.Warningf("Failed to delete Pod: %s", err)
		return model.RetVal{Msg: i18n.K8S.DeleteError, Attr: map[string]any{"Model": "Pod", "Error": err.Error()}}
	}
	return model.SuccessRetVal()
}

func DeletePodList(ctx context.Context, labels ...map[string]string) model.RetVal {
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
