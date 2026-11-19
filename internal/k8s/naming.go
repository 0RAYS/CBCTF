package k8s

import (
	"crypto/sha256"
	"fmt"
)

// victimResourceName keeps the name within the DNS label length even when the
// compose service or container key is long.
func victimResourceName(kind string, victimID uint, key string) string {
	hash := sha256.Sum256([]byte(key))
	return fmt.Sprintf("%s-%d-%x", kind, victimID, hash[:8])
}
