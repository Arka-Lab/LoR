package tools

import (
	"encoding/hex"

	"golang.org/x/crypto/sha3"
)

func SHA256(data string) []byte {
	h := sha3.New256()
	h.Write([]byte(data))
	return h.Sum(nil)
}

func SHA256str(data string) string {
	return hex.EncodeToString(SHA256(data))
}

func SHA256arr(data string) [8]int {
	var result [8]int
	for index, c := range SHA256(data) {
		result[index/8] ^= int(c) << (index % 8)
	}
	return result
}

func SHA256int(data string) int {
	var result int
	for _, c := range SHA256arr(data) {
		result ^= c
	}
	return result
}
