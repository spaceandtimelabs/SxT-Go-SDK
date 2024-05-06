package main

import (
	"crypto/ed25519"
	"os"
	"testing"

	"github.com/spaceandtimelabs/SxT-Go-SDK/utils"
)

var userId, privKeyB64, pubKeyB64 string
var privKey ed25519.PrivateKey
var pubKey ed25519.PublicKey

// Test Authentication
func TestAuthentication(t *testing.T){
	userId = os.Getenv("TEST_USER")
	privKeyB64 = os.Getenv("TEST_USER_PRIVKEY")
	pubKeyB64 = os.Getenv("TEST_USER_PUBKEY")

	_, _, privKeyBytes, pubKeyBytes, err := utils.Authenticate(userId, pubKeyB64, privKeyB64)

	if err != nil {
		t.Errorf("Autentication error %q", err)
	}

	pubKey = pubKeyBytes
	privKey = privKeyBytes

}

// Test SQL APIs
func TestSQLAPIs(t *testing.T){
	err := utils.SQLAPIs(privKey, pubKey)
	// Intentionally skip errors as we dont want to modify any records
	// This will always throw error
	if err == nil {
		t.Errorf("SQL API error %q", err)
	}
}



// Test Discovery APIs
func TestDiscoveryAPIs(t *testing.T){
	err := utils.DiscoveryAPIs()
	if err != nil {
		t.Errorf("Discovery APIs error %q", err)
	}
}