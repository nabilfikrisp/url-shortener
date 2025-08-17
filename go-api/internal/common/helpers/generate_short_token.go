package helpers

import (
	"crypto/sha1"
	"encoding/hex"
)

func GenerateShortToken(s string) string {
	hash := sha1.Sum([]byte(s))
	return hex.EncodeToString(hash[:])[:16]
}
