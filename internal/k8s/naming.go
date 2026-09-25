package k8s

import (
	"crypto/sha256"
	"fmt"
)

// Bound DNS label length even when compose service/container keys are long.
func victimResourceName(kind string, victimID uint, key string) string {
	hash := sha256.Sum256([]byte(key))
	return fmt.Sprintf("%s-%d-%x", kind, victimID, hash[:8])
}
