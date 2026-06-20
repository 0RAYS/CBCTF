package k8s

import (
	"CBCTF/internal/config"
	"CBCTF/internal/log"
	"fmt"

	"k8s.io/client-go/rest"

	netattclient "github.com/k8snetworkplumbingwg/network-attachment-definition-client/pkg/client/clientset/versioned"
	ovnclient "github.com/kubeovn/kube-ovn/pkg/client/clientset/versioned"
	"k8s.io/client-go/kubernetes"
	virtclient "kubevirt.io/client-go/kubevirt"
)

const (
	RoleLabel    = "role"
	ServiceLabel = "service"

	VictimPodTag    = "victim"
	GeneratorPodTag = "generator"
	FrpcPodTag      = "frpc"

	CaptureContainerName = "capture"
	NginxContainerName   = "nginx"
	FrpcContainerName    = "frpc"
)

var (
	kubeClient      *kubernetes.Clientset
	netattClient    *netattclient.Clientset
	ovnClient       *ovnclient.Clientset
	virtClient      *virtclient.Clientset
	kubeConfig      *rest.Config
	globalNamespace string
	nfsVolumeName   string
)

func Init() {
	globalNamespace = config.Env.K8S.Namespace
	nfsVolumeName = fmt.Sprintf("%s-shared-volume", globalNamespace)
	initClients()
	checkPermissions()
	checkResources()
}

func initClients() {
	var err error
	kubeConfig, err = rest.InClusterConfig()
	if err != nil {
		log.Logger.Fatalf("Failed to create in-cluster Kubernetes config: %s", err)
	}
	kubeConfig.QPS = 100
	kubeConfig.Burst = 150
	log.Logger.Info("Admin config loaded")
	kubeClient, err = kubernetes.NewForConfig(kubeConfig)
	if err != nil {
		log.Logger.Fatalf("Failed to init k8s client: %s", err)
	}
	netattClient, err = netattclient.NewForConfig(kubeConfig)
	if err != nil {
		log.Logger.Fatalf("Failed to init NetworkAttachmentDefinition client: %s", err)
	}
	ovnClient, err = ovnclient.NewForConfig(kubeConfig)
	if err != nil {
		log.Logger.Fatalf("Failed to init KubeOVN client: %s", err)
	}
	virtClient, err = virtclient.NewForConfig(kubeConfig)
	if err != nil {
		log.Logger.Fatalf("Failed to init KubeVirt client: %s", err)
	}
}
