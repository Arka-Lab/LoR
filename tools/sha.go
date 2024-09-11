package tools

import (
	"encoding/hex"

	"golang.org/x/crypto/sha3"
)

func SHA3_256(data string) string {
	h := sha3.New256()
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}
