package storage

import (
	"crypto/ed25519"
	"encoding/json"
	"io/fs"
	"log"
	"os"
)

type FileSessionStruct struct {
	AccessToken  string             `json:"accessToken"`
	RefreshToken string             `json:"refreshToken"`
	PrivateKey   ed25519.PrivateKey `json:"privateKey"`
	PublicKey    ed25519.PublicKey  `json:"publicKey"`
}

// Write session data to file: accessToken, refreshToken, publicKey, privateKey
// Note: This is for demo purposes only. To make a secure session, save the credentials in db or some other secure source
func FileWriteSession(userId, accessToken, refreshToken string, privateKey ed25519.PrivateKey, publicKey ed25519.PublicKey) (status bool) {
	if accessToken == "" || refreshToken == "" {
		return false
	}

	sessionStruct := FileSessionStruct{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		PublicKey:    publicKey,
		PrivateKey:   privateKey,
	}

	sessionData, err := json.Marshal(sessionStruct)

	if err != nil {
		log.Println(err.Error())
	} else {
		filepath := "./tmp/" + userId + ".txt"
		everyoneCanReadWriteAndExecute := fs.FileMode(0777)
		file, err := os.OpenFile(filepath, os.O_CREATE|os.O_WRONLY, everyoneCanReadWriteAndExecute)

		if err != nil {
			log.Println(err.Error())
			return false
		}

		defer file.Close()

		_, errWrite := file.Write(sessionData)
		if errWrite != nil {
			log.Println(errWrite.Error())
			return false
		}
	}

	return true
}

// Read session data from file: accessToken, refreshToken, publicKey, privateKey
// Note: This is for demo purposes only. To make a secure session, save the credentials in db or some other secure source
func FileReadSession(userId string) (sessionStruct FileSessionStruct, status bool) {
	content, err := os.ReadFile("./tmp/" + userId + ".txt")

	if err != nil || json.Unmarshal(content, &sessionStruct) != nil {
		return FileSessionStruct{}, false
	}

	return sessionStruct, true
}
