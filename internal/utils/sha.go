package utils

import "crypto/sha256"

func GetFirstTwoBytesFromSha256(v string) []byte {
	h := sha256.New()
	h.Write([]byte(v))
	sha256 := h.Sum(nil)
	return sha256[:2]
}
