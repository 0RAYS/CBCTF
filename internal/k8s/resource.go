package k8s

import (
	"CBCTF/internal/config"
	"CBCTF/internal/log"
	"context"
	"fmt"
	"time"
)

func initExternalNetworks() {
	if !config.Env.K8S.ExternalNetworks.Enabled {
		log.Logger.Warningf("External networks are not enabled, VPC will not work correctly")
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	interfaces := config.Env.K8S.ExternalNetworks.Interfaces
	externalNetworks = make([]ExternalNetwork, 0, len(interfaces))
	for _, iface := range interfaces {
		if iface.Interface == "" || iface.CIDR == "" || iface.Gateway == "" {
			continue
		}
		network := ExternalNetwork{
			SubnetName: fmt.Sprintf("%s-external-network-%s", globalNamespace, iface.Interface),
			Interface:  iface.Interface,
			CIDR:       iface.CIDR,
			Gateway:    iface.Gateway,
		}
		DeleteNetAttachDef(ctx, network.SubnetName)
		if _, ret := CreateNetAttachDef(ctx, CreateNetAttachDefOptions{
			Name:      network.SubnetName,
			Namespace: "kube-system",
			Config: fmt.Sprintf(`{
				"cniVersion": "0.3.0",
				"type": "macvlan",
				"master": %q,
				"mode": "bridge"
			}`, network.Interface),
		}); !ret.OK {
			continue
		}

		DeleteSubnet(ctx, network.SubnetName)
		if _, ret := CreateSubnet(ctx, CreateSubnetOptions{
			Name:     network.SubnetName,
			CIDR:     network.CIDR,
			Gateway:  network.Gateway,
			Provider: fmt.Sprintf("%s.kube-system", network.SubnetName),
		}); !ret.OK {
			continue
		}
		log.Logger.Infof("External network initialized: subnet=%s interface=%s", network.SubnetName, network.Interface)
		externalNetworks = append(externalNetworks, network)
	}
}

func checkResources() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	if _, ret := GetNamespace(ctx, globalNamespace); !ret.OK {
		log.Logger.Fatalf("Namespace %q not found: %s", globalNamespace, ret.Msg)
	}
	if _, ret := GetPVC(ctx, nfsVolumeName); !ret.OK {
		log.Logger.Warningf("PersistentVolumeClaim %q not found: %s", nfsVolumeName, ret.Msg)
		log.Logger.Warningf("Dynamic attachments will not be generated correctly")
	}
	for _, network := range externalNetworks {
		if _, ret := GetNetAttachDef(ctx, network.SubnetName, "kube-system"); !ret.OK {
			log.Logger.Warningf("Network %q not found: %s", network.SubnetName, ret.Msg)
		}
		if _, ret := GetSubnet(ctx, network.SubnetName); !ret.OK {
			log.Logger.Warningf("Subnet %q not found: %s", network.SubnetName, ret.Msg)
		}
	}
	log.Logger.Info("K8s resource check completed")
}
