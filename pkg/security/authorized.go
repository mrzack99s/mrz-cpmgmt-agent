package security

import (
	"crypto/sha256"
	"fmt"

	encryptions "github.com/mrzack99s/mrz-cpmgmt-agent/pkg/encryptions/aes"
)

func CheckAuthorized(reqAuthorized string) bool {

	hash := sha256.Sum256([]byte(encryptions.AES_KEY))
	authKey := fmt.Sprintf("%x", hash[:])
	if reqAuthorized == authKey {
		return true
	}

	return false
}
