package k8s

import (
	"CBCTF/internal/config"
	"CBCTF/internal/log"
	"fmt"
	"os"

	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	netattclient "github.com/k8snetworkplumbingwg/network-attachment-definition-client/pkg/client/clientset/versioned"
	ovnclient "github.com/kubeovn/kube-ovn/pkg/client/clientset/versioned"
	"k8s.io/client-go/kubernetes"
)

const (
	VictimPodTag    = "victim"
	GeneratorPodTag = "generator"
	FrpcPodTag      = "frpc"
)

var (
	kubeClient       *kubernetes.Clientset
	netattClient     *netattclient.Clientset
	ovnClient        *ovnclient.Clientset
	kubeConfig       *rest.Config
	globalNamespace  string
	externalNetworks []ExternalNetwork
	nfsVolumeName    string
)

type ExternalNetwork struct {
	SubnetName string
	Interface  string
	CIDR       string
	Gateway    string
}

func Init() {
	globalNamespace = config.Env.K8S.Namespace
	nfsVolumeName = fmt.Sprintf("%s-shared-volume", globalNamespace)
	initClients()
	checkPermissions()
	initExternalNetworks()
	checkResources()
}

func initClients() {
	var err error
	kubeConfig, err = rest.InClusterConfig()
	if err != nil {
		if _, err = os.Stat(config.Env.K8S.Config); err != nil {
			log.Logger.Fatalf("Invalid config.k8s.config.admin: %s", err)
		}
		kubeConfig, err = clientcmd.BuildConfigFromFlags("", config.Env.K8S.Config)
		if err != nil {
			log.Logger.Fatalf("Failed to create client config: %s", err)
		}
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
}
