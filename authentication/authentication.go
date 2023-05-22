package authentication

import (
	"bytes"
	"crypto"
	"crypto/ed25519"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sxt-sdks/helpers"

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

	apiEndPoint, _ := helpers.ReadEndPoint()
	codeEndPoint := apiEndPoint + "/auth/code"

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
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	sb := string(body)
	return sb
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

	apiEndPoint, _ := helpers.ReadEndPoint()
	scheme, _ := helpers.ReadScheme()
	tokenEndPoint := apiEndPoint + "/auth/token"

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
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	token = string(body)

	return token
}

// Get new accesstoken and refreshToken from provided `refreshToken`
func RefreshToken(refreshToken string) (tokenStruct TokenStruct, status bool) {
	apiEndPoint, _ := helpers.ReadEndPoint()
	tokenEndPoint := apiEndPoint + "/auth/refresh"

	client := http.Client{}
	req, err := http.NewRequest("POST", tokenEndPoint, nil)
	if err != nil {
		return TokenStruct{}, false
	}

	req.Header.Add("Authorization", "Bearer "+refreshToken)

	res, err := client.Do(req)
	if err != nil {
		return TokenStruct{}, false
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return TokenStruct{}, false
	}

	e := json.Unmarshal(body, &tokenStruct)
	if e != nil {
		return TokenStruct{}, false
	}

	return tokenStruct, true
}

// validate access token, if its active
func ValidateToken(accessToken string) (status bool) {
	apiEndPoint, _ := helpers.ReadEndPoint()
	tokenEndPoint := apiEndPoint + "/auth/validtoken"

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
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatalln(err)
	}

	strLen := len(string(body))
	if strLen <= 0 {
		status = false
	} else {
		status = true
	}

	return status
}

// Logout user
func Logout() {

	apiEndPoint, _ := helpers.ReadEndPoint()
	tokenEndPoint := apiEndPoint + "/auth/logout"

	client := http.Client{}
	req, err := http.NewRequest("POST", tokenEndPoint, nil)
	if err != nil {
		log.Fatalln(err)
	}

	req.Header.Add("Authorization", "Bearer "+os.Getenv("accessToken"))

	client.Do(req)
}
