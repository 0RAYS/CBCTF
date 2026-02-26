package k8s

import (
	"CBCTF/internal/log"
	"context"
	"time"
)

func checkResources() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	checkNamespace(ctx)
	checkNFSVolume(ctx)
	checkExternalNetwork(ctx)
	log.Logger.Info("K8s resource check completed")
}

func checkNamespace(ctx context.Context) {
	if _, ret := GetNamespace(ctx, globalNamespace); !ret.OK {
		log.Logger.Fatalf("Namespace %q not found: %s", globalNamespace, ret.Msg)
	}
}

func checkExternalNetwork(ctx context.Context) {
	if _, ret := GetSubnet(ctx, externalSubnetName); !ret.OK {
		log.Logger.Warningf("Subnet %q not found: %s", externalSubnetName, ret.Msg)
		log.Logger.Warningf("VPC network will not work correctly")
	}

	if _, ret := GetNetAttachDef(ctx, externalSubnetName, "kube-system"); !ret.OK {
		log.Logger.Warningf("NetworkAttachmentDefinition %q not found in kube-system: %s", externalSubnetName, ret.Msg)
		log.Logger.Warningf("VPC network will not work correctly")
	}
}

func checkNFSVolume(ctx context.Context) {
	if _, ret := GetPVC(ctx, nfsVolumeName); !ret.OK {
		log.Logger.Warningf("PersistentVolumeClaim %q not found: %s", nfsVolumeName, ret.Msg)
		log.Logger.Warningf("Dynamic attachments will not be generated correctly")
	}
}
