package authentication

import (
	"bytes"
	"crypto"
	"crypto/ed25519"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/spaceandtimelabs/SxT-Go-SDK/helpers"

	b64 "encoding/base64"
)

type AuthCodeStruct struct {
	AuthCode string `json:"authCode"`
}

type TokenStruct struct {
	AccessToken         string `json:"accessToken"`
	RefreshToken        string `json:"refreshToken"`
	AccessTokenExpires  int    `json:"accessTokenExpires"`
	RefreshTokenExpires int    `json:"refreshTokenExpires"`
}

// Generate auth code
func GenerateAuthCode(userId, joinCode string) (authCode string) {
	codeEndPoint := helpers.GetAuthenticationEndpoint("code")
	postBody, _ := json.Marshal(map[string]string{
		"userId":   userId,
		"joinCode": joinCode,
	})
	responseBody := bytes.NewBuffer(postBody)
	resp, err := http.Post(codeEndPoint, "application/json", responseBody)

	//Handle Error
	if err != nil {
		log.Fatalf("An Error Occured %v", err)
	}
	defer resp.Body.Close()

	//Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	return string(body)
}

// Generate Encoded signature and base64 public key
func GenerateKeys(authCode string, pubkey ed25519.PublicKey, privkey ed25519.PrivateKey) (encodedSignature, base64PublicKey string) {
	signature, eSignature := privkey.Sign(rand.Reader, []byte(authCode), crypto.Hash(0))
	if eSignature != nil {
		log.Fatalln("Signature Error", eSignature.Error())
	} else {
		encodedSignature = hex.EncodeToString(signature)

		base64PublicKey = b64.StdEncoding.EncodeToString([]byte(pubkey))
	}

	return encodedSignature, base64PublicKey
}

// Generate accessToken, refreshToken
func GenerateToken(userId, authCode, encodedSignature, base64PublicKey string) (token string) {
	tokenEndPoint := helpers.GetAuthenticationEndpoint("token")
	scheme, _ := helpers.ReadScheme()
	postBody, _ := json.Marshal(map[string]string{
		"userId":    userId,
		"authCode":  authCode,
		"key":       base64PublicKey,
		"signature": encodedSignature,
		"scheme":    scheme,
	})
	responseBody := bytes.NewBuffer(postBody)
	resp, err := http.Post(tokenEndPoint, "application/json", responseBody)

	//Handle Error
	if err != nil {
		log.Fatalf("An Error Occured %v", err)
	}
	defer resp.Body.Close()

	//Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	token = string(body)

	return token
}

// Get new accesstoken and refreshToken from provided `refreshToken`
func RefreshToken(refreshToken string) (tokenStruct TokenStruct, status bool) {
	tokenEndPoint := helpers.GetAuthenticationEndpoint("refresh")
	req, err := http.NewRequest("POST", tokenEndPoint, nil)
	if err != nil {
		return TokenStruct{}, false
	}

	req.Header.Add("Authorization", "Bearer "+refreshToken)

	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return TokenStruct{}, false
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return TokenStruct{}, false
	}

	err = json.Unmarshal(body, &tokenStruct)
	if err != nil {
		return TokenStruct{}, false
	}

	return tokenStruct, true
}

// validate access token, if its active
func ValidateToken(accessToken string) (status bool) {
	tokenEndPoint := helpers.GetAuthenticationEndpoint("validtoken")
	client := http.Client{}
	req, err := http.NewRequest("GET", tokenEndPoint, nil)
	if err != nil {
		log.Fatalln(err)
	}

	req.Header.Add("Authorization", "Bearer "+accessToken)

	res, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatalln(err)
	}

	strLen := len(string(body))

	return strLen > 0
}

// Logout user
func Logout() {
	tokenEndPoint := helpers.GetAuthenticationEndpoint("logout")
	client := http.Client{}
	req, err := http.NewRequest("POST", tokenEndPoint, nil)
	if err != nil {
		log.Fatalln(err)
	}

	req.Header.Add("Authorization", "Bearer "+os.Getenv("accessToken"))

	client.Do(req)
}
