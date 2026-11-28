package k8s

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"maps"
	"strings"
	"time"

	corev1 "k8s.io/api/core/v1"
	apierror "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labelselector "k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/wait"

	"CBCTF/internal/config"
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
)

type Network struct {
	Interface    string
	IPv4         string
	MAC          string
	CIDR         string
	Subnet       string
	NetAttachDef string
	Routes       []model.Route
	// GuestAddress 表示地址配在虚拟机内部,virt-launcher 的附加网卡只保留二层连接。
	GuestAddress bool
}

// ovnAnnotations 生成 Multus 与 Kube-OVN 的网卡注解。
// 静态地址与静态路由都通过 pod 注解交给 Kube-OVN CNI。
// default_route 固定为 false:工作负载的默认路由属于集群网络,靶场网段只做附加网卡,
// 跨网段可达性由静态路由给出。这样 kube-ovn 无论作主 CNI 还是副 CNI,行为都一致,
// 也不会把工作负载从集群网络上摘下来。
func ovnAnnotations(networks []Network) map[string]string {
	annotations := make(map[string]string)
	for _, network := range networks {
		annotations["k8s.v1.cni.cncf.io/networks"] += fmt.Sprintf(",%s/%s", globalNamespace, network.NetAttachDef)
		annotations["k8s.v1.cni.cncf.io/networks"] = strings.Trim(annotations["k8s.v1.cni.cncf.io/networks"], ",")
		prefix := fmt.Sprintf("%s.%s.ovn.kubernetes.io/", network.NetAttachDef, globalNamespace)
		annotations[prefix+"logical_switch"] = network.Subnet
		if network.MAC != "" {
			annotations[prefix+"mac_address"] = network.MAC
		}
		annotations[prefix+"default_route"] = "false"
		if network.GuestAddress {
			// 虚拟机的地址由 guest 自己持有,Pod 侧再放一份会让 Pod 内核
			// 直接吃掉发往该地址的流量,虚拟机反而收不到。
			continue
		}
		annotations[prefix+"ip_address"] = network.IPv4
		if len(network.Routes) > 0 {
			if data, err := json.Marshal(network.Routes); err == nil {
				annotations[prefix+"routes"] = string(data)
			}
		}
	}
	return annotations
}

type CreatePodOptions struct {
	SubmitOnly         bool
	Name               string
	Labels             map[string]string
	Annotations        map[string]string
	Networks           []Network
	Containers         []corev1.Container
	InitContainers     []corev1.Container
	Volumes            []corev1.Volume
	ServiceAccountName string
}

// runtimeServiceAccount is the ServiceAccount of the runtime workloads. It is
// empty when the platform runs with the ServiceAccount of the namespace.
func runtimeServiceAccount() string {
	return config.Env.K8S.ServiceAccount
}

func CreatePod(ctx context.Context, options CreatePodOptions) (*corev1.Pod, model.RetVal) {
	affinity, err := imageFailureAffinity(ctx, append(containerImages(options.Containers), containerImages(options.InitContainers)...))
	if err != nil {
		return nil, model.FailRetVal(i18n.K8S.CreateError, "Pod", options.Name, err)
	}
	annotations := make(map[string]string)
	maps.Copy(annotations, options.Annotations)
	maps.Copy(annotations, ovnAnnotations(options.Networks))
	if len(annotations) == 0 {
		annotations = nil
	}
	pod := &corev1.Pod{
		OwnerReferences: resourceOwners(ctx),
		Name:            options.Name,
		Namespace:       globalNamespace,
		Labels:          options.Labels,
		Annotations:     annotations,
		Spec: corev1.PodSpec{
			Affinity:                      affinity,
			EnableServiceLinks:            new(false),
			AutomountServiceAccountToken:  new(false),
			Containers:                    options.Containers,
			InitContainers:                options.InitContainers,
			Volumes:                       options.Volumes,
			ServiceAccountName:            options.ServiceAccountName,
			TerminationGracePeriodSeconds: new(int64(3)),
			RestartPolicy:                 corev1.RestartPolicyNever,
		},
	}
	pod, err = kubeClient.CoreV1().Pods(globalNamespace).Create(ctx, pod, metav1.CreateOptions{})
	if err != nil {
		log.Logger.Warningf("Failed to create Pod %s: %s", options.Name, err)
		return nil, model.FailRetVal(i18n.K8S.CreateError, "Pod", options.Name, err)
	}
	if options.SubmitOnly {
		return pod, model.SuccessRetVal()
	}
	ready, err := waitPodReady(ctx, pod)
	if err != nil {
		return nil, model.FailRetVal(i18n.K8S.PodRunError, "Pod", pod.Name, err).WithAttrs(map[string]any{
			"Pod": pod.Name, "Phase": pod.Status.Phase, "Reason": err.Error(),
		})
	}
	return ready, model.SuccessRetVal()
}

func GetPod(ctx context.Context, name string) (*corev1.Pod, model.RetVal) {
	pod, err := kubeClient.CoreV1().Pods(globalNamespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, getFailed("Pod", name, err)
	}
	return pod, model.SuccessRetVal()
}

func ListPods(ctx context.Context, labels ...map[string]string) (*corev1.PodList, model.RetVal) {
	c, err := cachedObjects(ctx, "pods")
	if err != nil {
		return nil, model.FailRetVal(i18n.K8S.GetError, "Pod", "", err)
	}
	selector := labelselector.Everything()
	if len(labels) > 0 {
		selector = labelselector.SelectorFromSet(labels[0])
	}
	list := &corev1.PodList{}
	for _, obj := range c.informer.GetStore().List() {
		pod := obj.(*corev1.Pod)
		if selector.Matches(labelselector.Set(pod.Labels)) {
			list.Items = append(list.Items, *pod.DeepCopy())
		}
	}
	return list, model.SuccessRetVal()
}

func listPodsDirect(ctx context.Context, labels ...map[string]string) (*corev1.PodList, model.RetVal) {
	var options metav1.ListOptions
	if len(labels) > 0 {
		options.LabelSelector = labelselector.SelectorFromSet(labels[0]).String()
	}
	podList, err := kubeClient.CoreV1().Pods(globalNamespace).List(ctx, options)
	if err != nil {
		log.Logger.Warningf("Failed to list Pods (%s): %s", options.LabelSelector, err)
		return nil, model.FailRetVal(i18n.K8S.GetError, "Pod", options.LabelSelector, err)
	}
	return podList, model.SuccessRetVal()
}

const MaxPodLogLines int64 = 10000
const MaxPodLogBytes int64 = 1 << 20

// networkStatusAnnotation 是 Multus 汇总各网卡地址的注解。
const networkStatusAnnotation = "k8s.v1.cni.cncf.io/network-status"

type networkStatus struct {
	Name      string   `json:"name"`
	Interface string   `json:"interface"`
	IPs       []string `json:"ips"`
	MAC       string   `json:"mac"`
	Default   bool     `json:"default"`
}

// podAddresses 返回 Pod 各网卡上的地址。
// 注解缺失时无法判断实际地址,返回 false 让调用方跳过校验。
func podAddresses(pod *corev1.Pod) (map[string]struct{}, bool) {
	data, ok := pod.Annotations[networkStatusAnnotation]
	if !ok {
		return nil, false
	}
	var statuses []networkStatus
	if err := json.Unmarshal([]byte(data), &statuses); err != nil {
		return nil, false
	}
	addresses := make(map[string]struct{})
	for _, status := range statuses {
		for _, ip := range status.IPs {
			addresses[ip] = struct{}{}
		}
	}
	return addresses, true
}

// verifyPodAddresses 校验 Pod 实际拿到的地址与拓扑规划一致。
// 注解没被 CNI 采纳时 Pod 会拿到别的地址,靶场拓扑随即失效,必须在就绪前发现。
func verifyPodAddresses(pod *corev1.Pod, networks model.Networks) error {
	addresses, ok := podAddresses(pod)
	if !ok {
		return nil
	}
	for _, network := range networks {
		expected := network.Attachment.IP
		if expected == "" {
			continue
		}
		if _, found := addresses[expected]; !found {
			return fmt.Errorf("pod %s got addresses %v, expected %s on %s", pod.Name, addresses, expected, network.Attachment.Name)
		}
	}
	return nil
}

func podLogOptions(containerName string, lines int64) *corev1.PodLogOptions {
	if lines <= 0 {
		lines = 1000
	}
	lines = min(lines, MaxPodLogLines)
	return &corev1.PodLogOptions{Container: containerName, TailLines: &lines, LimitBytes: new(MaxPodLogBytes)}
}

func GetPodLogs(ctx context.Context, podName, containerName string, lines int64) (string, model.RetVal) {
	options := podLogOptions(containerName, lines)
	target := fmt.Sprintf("%s/%s", podName, containerName)
	podLogs, err := kubeClient.CoreV1().Pods(globalNamespace).GetLogs(podName, options).Stream(ctx)
	if err != nil {
		log.Logger.Warningf("Failed to get Pod Logs %s: %s", target, err)
		return "", model.FailRetVal(i18n.K8S.GetError, "PodLog", target, err)
	}
	defer func(podLogs io.ReadCloser) {
		_ = podLogs.Close()
	}(podLogs)
	buf, err := io.ReadAll(io.LimitReader(podLogs, MaxPodLogBytes))
	if err != nil {
		log.Logger.Warningf("Failed to read Pod Logs %s: %s", target, err)
		return "", model.FailRetVal(i18n.K8S.GetError, "PodLog", target, err)
	}
	return string(buf), model.SuccessRetVal()
}

// DeletePod checks the observed UID to avoid deleting a replacement Pod.
func DeletePod(ctx context.Context, name string, uid types.UID) model.RetVal {
	if name == "" || uid == "" {
		return model.FailRetVal(i18n.K8S.DeleteError, "Pod", name, errMissingPodIdentity)
	}
	options := metav1.DeleteOptions{Preconditions: &metav1.Preconditions{UID: &uid}}
	err := kubeClient.CoreV1().Pods(globalNamespace).Delete(ctx, name, options)
	if err != nil && !apierror.IsNotFound(err) {
		log.Logger.Warningf("Failed to delete Pod %s: %s", name, err)
		return model.FailRetVal(i18n.K8S.DeleteError, "Pod", name, err)
	}
	return model.SuccessRetVal()
}

func DeletePodCollection(ctx context.Context, labels map[string]string) model.RetVal {
	return deleteCollection("Pod", "Pods", labels, func(options metav1.ListOptions) error {
		return kubeClient.CoreV1().Pods(globalNamespace).DeleteCollection(ctx, metav1.DeleteOptions{}, options)
	})
}

// DeletePodAndWait waits for the deletion to complete: acceptance is not
// completion, PVC mounts and exec may still be active.
func DeletePodAndWait(ctx context.Context, name string, uid types.UID) model.RetVal {
	if ret := DeletePod(ctx, name, uid); !ret.OK {
		return ret
	}
	err := wait.PollUntilContextCancel(ctx, 500*time.Millisecond, true, func(ctx context.Context) (bool, error) {
		pod, err := cachedPod(ctx, name)
		if err == nil && pod == nil {
			_, err = kubeClient.CoreV1().Pods(globalNamespace).Get(ctx, name, metav1.GetOptions{})
			return apierror.IsNotFound(err), errUnlessNotFound(err)
		}
		if err == nil && pod != nil && pod.UID != uid {
			return false, fmt.Errorf("pod %s was replaced during cleanup", name)
		}
		return false, err
	})
	if err != nil {
		return model.FailRetVal(i18n.K8S.DeleteError, "Pod", name, err)
	}
	return model.SuccessRetVal()
}
