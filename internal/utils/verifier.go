package utils

import (
	"crypto/sha256"
	"crypto/subtle"
	"fmt"
)

// HashVerifier hashes one-time verification secrets before embedding them in JWTs.
func HashVerifier(verifier string) string {
	inner := sha256.Sum256([]byte(verifier))
	outer := sha256.Sum256([]byte(verifier + fmt.Sprintf("%x", inner)))
	return fmt.Sprintf("%x", outer)
}

func CompareVerifier(verifier string, hash string) bool {
	return subtle.ConstantTimeCompare([]byte(HashVerifier(verifier)), []byte(hash)) == 1
}
