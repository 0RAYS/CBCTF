package k8s

import (
	"CBCTF/internal/config"
	"CBCTF/internal/log"
	"context"
	"fmt"
	"os"
	"time"

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
	nfsVolumeName = fmt.Sprintf("%s-nfs-volume", globalNamespace)
	initClients()
	checkPermissions()
}

func InitResources() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	CreateNamespace(ctx, CreateNamespaceOptions{Name: globalNamespace})
	initNFSVolume(ctx)
	initExternalNetwork(ctx)

	os.Exit(0)
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

func initExternalNetwork(ctx context.Context) {
	if _, ret := GetSubnet(ctx, externalSubnetName); !ret.OK {
		if _, ret := CreateSubnet(ctx, CreateSubnetOptions{
			Name:       externalSubnetName,
			CIDR:       config.Env.K8S.ExternalNetwork.CIDR,
			Gateway:    config.Env.K8S.ExternalNetwork.Gateway,
			ExcludeIPs: config.Env.K8S.ExternalNetwork.ExcludeIPs,
			Provider:   fmt.Sprintf("%s.kube-system", externalSubnetName),
		}); !ret.OK {
			log.Logger.Fatal("Failed to init external network")
		}
	} else {
		log.Logger.Info("ExternalNetworkSubnet is already exists")
	}
	if _, ret := GetNetAttachDef(ctx, externalSubnetName, "kube-system"); !ret.OK {
		if _, ret := CreateNetAttachDef(ctx, CreateNetAttachDefOptions{
			Name:      externalSubnetName,
			Namespace: "kube-system",
			Config: fmt.Sprintf(`{
			"cniVersion": "0.3.0",
			"type": "macvlan",
			"master": "%s",
			"mode": "bridge",
			"ipam": {
				"type": "kube-ovn",
				"server_socket": "/run/openvswitch/kube-ovn-daemon.sock",
				"provider": "%s.kube-system"
			}
		}`, config.Env.K8S.ExternalNetwork.Interface, externalSubnetName),
		}); !ret.OK {
			log.Logger.Fatal("Failed to init external network attachment definition")
		}
	} else {
		log.Logger.Info("ExternalNetworkAttachDef is already exists")
	}
}

func initNFSVolume(ctx context.Context) {
	if _, ret := GetPV(ctx, nfsVolumeName); !ret.OK {
		if _, ret := CreatePV(ctx, CreatePVOptions{
			Name:    nfsVolumeName,
			Server:  config.Env.NFS.Server,
			Path:    config.Env.NFS.Path,
			Storage: config.Env.NFS.Storage,
		}); !ret.OK {
			log.Logger.Fatal("Failed to init pv")
		}
	} else {
		log.Logger.Info("NFS Volume PV is already exists")
	}
	if _, ret := GetPVC(ctx, nfsVolumeName); !ret.OK {
		if _, ret := CreatePVC(ctx, CreatePVCOptions{
			Name:    nfsVolumeName,
			Storage: config.Env.NFS.Storage,
		}); !ret.OK {
			log.Logger.Fatalf("Failed to init pvc")
		}
	} else {
		log.Logger.Info("NFS Volume VPC is already exists")
	}
	log.Logger.Infof("Please mount the nfs server %s:%s at %s manually", config.Env.NFS.Server, config.Env.NFS.Path, config.Env.Path)
}
