package k8s

import (
	"CBCTF/internal/config"
	"CBCTF/internal/log"
	"fmt"
	"os"

	"k8s.io/client-go/rest"

	kubeovnclient "github.com/JBNRZ/kubeovn-api/pkg/client/clientset/versioned"
	netattclient "github.com/k8snetworkplumbingwg/network-attachment-definition-client/pkg/client/clientset/versioned"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	VictimPodTag    = "victim"
	GeneratorPodTag = "generator"
	FrpcPodTag      = "frpc"
)

const (
	VPCNetworkTolerationKey = "vpc-network"
	VPCNetworkTolerationVal = "unacceptable"
)

var (
	kubeClient         *kubernetes.Clientset
	netattClient       *netattclient.Clientset
	kubeOVNClient      *kubeovnclient.Clientset
	kubeConfig         *rest.Config
	globalNamespace    string
	externalSubnetName string
	nfsVolumeName      string
)

func Init() {
	globalNamespace = config.Env.K8S.Namespace
	externalSubnetName = fmt.Sprintf("%s-external-network", globalNamespace)
	nfsVolumeName = fmt.Sprintf("%s-shared-volume", globalNamespace)
	initClients()
	checkResources()
	checkPermissions()
}

func initClients() {
	if _, err := os.Stat(config.Env.K8S.Config); err != nil {
		log.Logger.Fatalf("Invalid config.k8s.config.admin: %s", err)
	}
	adminAPIConfig, err := clientcmd.LoadFromFile(config.Env.K8S.Config)
	if err != nil {
		log.Logger.Fatalf("Failed to load admin config: %s", err)
	}
	kubeConfig, err = clientcmd.NewNonInteractiveClientConfig(*adminAPIConfig, adminAPIConfig.CurrentContext, &clientcmd.ConfigOverrides{}, nil).ClientConfig()
	if err != nil {
		log.Logger.Fatalf("Failed to create client config: %s", err)
	}
	kubeConfig.QPS = 700
	kubeConfig.Burst = 1000
	log.Logger.Info("Admin config loaded")
	kubeClient, err = kubernetes.NewForConfig(kubeConfig)
	if err != nil {
		log.Logger.Fatalf("Failed to init k8s client: %s", err)
	}
	netattClient, err = netattclient.NewForConfig(kubeConfig)
	if err != nil {
		log.Logger.Fatalf("Failed to init NetworkAttachmentDefinition client: %s", err)
	}
	kubeOVNClient, err = kubeovnclient.NewForConfig(kubeConfig)
	if err != nil {
		log.Logger.Fatalf("Failed to init KubeOVN client: %s", err)
	}
}
