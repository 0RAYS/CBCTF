package k8s

import (
	"fmt"

	netattclient "github.com/k8snetworkplumbingwg/network-attachment-definition-client/pkg/client/clientset/versioned"
	ovnclient "github.com/kubeovn/kube-ovn/pkg/client/clientset/versioned"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/util/flowcontrol"
	virtclient "kubevirt.io/client-go/kubevirt"

	"CBCTF/internal/config"
	"CBCTF/internal/log"
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
	kubeClient      kubernetes.Interface
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
	configureClientRateLimit(kubeConfig)
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

func configureClientRateLimit(config *rest.Config) {
	config.QPS = 100
	config.Burst = 150
	// Share one process budget across core, Multus, Kube-OVN and KubeVirt clients.
	// Otherwise each independently creates a limiter, multiplying API-server bursts.
	config.RateLimiter = flowcontrol.NewTokenBucketRateLimiter(config.QPS, config.Burst)
	rest.AddUserAgent(config, "cbctf")
}
