package helpers

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

// func create Keys : ed25519 public and private keys
func CreateKey() (ed25519.PublicKey, ed25519.PrivateKey) {
	rng := rand.Reader
	publicRoot, privateRoot, _ := ed25519.GenerateKey(rng)

	fmt.Println("Public", base64.StdEncoding.EncodeToString(publicRoot))
	fmt.Println("Private", base64.StdEncoding.EncodeToString(privateRoot[0:32]))
	return publicRoot, privateRoot
}
