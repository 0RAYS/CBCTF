package k8s

import (
	"CBCTF/internal/log"
	"context"
	"time"
)

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
	log.Logger.Info("K8s resource check completed")
}
