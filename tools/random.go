package tools

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
)

func RandomSet(data []string, k int) ([]string, []string) {
	arr := make([]string, len(data))
	rnd := make([]int, 0)
	copy(arr, data)

	for i := 0; i < k; i++ {
		if len(rnd) == 0 {
			rnd = SHA256Arr(arr[i:])
		}
		index := rnd[0] % (len(arr) - i)
		arr[i], arr[index], rnd = arr[index], arr[i], rnd[1:]
	}
	return arr[:k], arr[k:]
}

func GeneratePrivateKey(size int) (*rsa.PrivateKey, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, size)
	if err != nil {
		return nil, err
	}
	if err := privateKey.Validate(); err != nil {
		return nil, err
	}
	return privateKey, nil
}

func SignWithPrivateKey(data []byte, privateKey *rsa.PrivateKey) ([]byte, error) {
	hashed := sha256.Sum256(data)
	return rsa.SignPSS(rand.Reader, privateKey, crypto.SHA256, hashed[:], nil)
}

func VerifyWithPublicKey(data []byte, signature []byte, publicKey *rsa.PublicKey) error {
	hashed := sha256.Sum256(data)
	return rsa.VerifyPSS(publicKey, crypto.SHA256, hashed[:], signature, nil)
}

func SignWithPrivateKeyStr(data string, privateKey *rsa.PrivateKey) (string, error) {
	signature, err := SignWithPrivateKey([]byte(data), privateKey)
	if err != nil {
		return "", err
	}
	return string(signature), nil
}

func VerifyWithPublicKeyStr(data string, signature string, publicKey *rsa.PublicKey) error {
	return VerifyWithPublicKey([]byte(data), []byte(signature), publicKey)
}
