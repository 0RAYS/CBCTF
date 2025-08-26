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

var (
	kubeClient         *kubernetes.Clientset
	natattClient       *netattclient.Clientset
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
}

func InitResources() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	updateNodeIPs(ctx)
	CreateNamespace(ctx, CreateNamespaceOptions{Name: globalNamespace})
	initNFSVolume(ctx)
	initExternalNetwork(ctx)

	if err := config.Save(config.Env); err != nil {
		log.Logger.Fatalf("Failed to update config: %s", err)
	}
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
	natattClient, err = netattclient.NewForConfig(kubeConfig)
	if err != nil {
		log.Logger.Fatalf("Failed to init NetworkAttachmentDefinition client: %s", err)
	}
	kubeOVNClient, err = kubeovnclient.NewForConfig(kubeConfig)
	if err != nil {
		log.Logger.Fatalf("Failed to init KubeOVN client: %s", err)
	}
}

func updateNodeIPs(ctx context.Context) {
	ips, ok, _ := GetNodeIPList(ctx)
	if !ok {
		log.Logger.Fatalf("Failed to get node IP list")
	}
	for _, ip := range ips {
		config.Env.K8S.Nodes = append(config.Env.K8S.Nodes, struct {
			IP     string `mapstructure:"ip" json:"ip" msgpack:"ip"`
			Public bool   `mapstructure:"public" json:"public" msgpack:"public"`
		}{IP: ip, Public: true})
	}
}

func initExternalNetwork(ctx context.Context) {
	if _, ok, _ := GetSubnet(ctx, externalSubnetName); !ok {
		if _, ok, _ := CreateSubnet(ctx, CreateSubnetOptions{
			Name:       externalSubnetName,
			CIDR:       config.Env.K8S.ExternalNetwork.CIDR,
			Gateway:    config.Env.K8S.ExternalNetwork.Gateway,
			ExcludeIPs: config.Env.K8S.ExternalNetwork.ExcludeIPs,
			Provider:   fmt.Sprintf("%s.kube-system", externalSubnetName),
		}); !ok {
			log.Logger.Fatal("Failed to init external network")
		}
	} else {
		log.Logger.Info("ExternalNetworkSubnet is already exists")
	}
	if _, ok, _ := GetNetAttachDef(ctx, externalSubnetName, "kube-system"); !ok {
		if _, ok, _ := CreateNetAttachDef(ctx, CreateNetAttachDefOptions{
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
		}); !ok {
			log.Logger.Fatal("Failed to init external network attachment definition")
		}
	} else {
		log.Logger.Info("ExternalNetworkAttachDef is already exists")
	}
}

func initNFSVolume(ctx context.Context) {
	if _, ok, _ := GetPV(ctx, nfsVolumeName); !ok {
		if _, ok, _ := CreatePV(ctx, CreatePVOptions{
			Name:    nfsVolumeName,
			Server:  config.Env.NFS.Server,
			Path:    config.Env.NFS.Path,
			Storage: config.Env.NFS.Storage,
		}); !ok {
			log.Logger.Fatal("Failed to init pv")
		}
	} else {
		log.Logger.Info("NFS Volume PV is already exists")
	}
	if _, ok, _ := GetPVC(ctx, nfsVolumeName); !ok {
		if _, ok, _ := CreatePVC(ctx, CreatePVCOptions{
			Name:    nfsVolumeName,
			Storage: config.Env.NFS.Storage,
		}); !ok {
			log.Logger.Fatalf("Failed to init pvc")
		}
	} else {
		log.Logger.Info("NFS Volume VPC is already exists")
	}
	log.Logger.Infof("Please mount the nfs server %s:%s at %s manually", config.Env.NFS.Server, config.Env.NFS.Path, config.Env.Path)
}
