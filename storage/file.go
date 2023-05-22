package storage

import (
	"crypto/ed25519"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

type FileSessionStruct struct{
	AccessToken string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	PrivateKey ed25519.PrivateKey `json:"privateKey"`
	PublicKey ed25519.PublicKey `json:"publicKey"`
}

// Write session data to file: accessToken, refreshToken, publicKey, privateKey
// Note: This is for demo purposes only. To make a secure session, save the credentials in db or some other secure source
func FileWriteSession(userId, accessToken, refreshToken string, privateKey ed25519.PrivateKey, publicKey ed25519.PublicKey)(status bool){

	if accessToken == "" || refreshToken == "" {
		return false
	}

 	sessionStruct := FileSessionStruct{
		AccessToken: accessToken,
		RefreshToken: refreshToken,
		PublicKey: publicKey,
		PrivateKey: privateKey,
	}

	sessionData, e := json.Marshal(sessionStruct)

	if e != nil {
		log.Println(e.Error())
	} else {
		

		f, errCreate := os.OpenFile("./tmp/" + userId + ".txt", os.O_CREATE|os.O_WRONLY, 0777)
		if errCreate != nil {
			log.Println(errCreate.Error())
			return false
		}
	
		defer f.Close()
	
		_, errWrite := f.Write(sessionData)
		if errWrite != nil {
			log.Println(errWrite.Error())
			return false
		}
	}

	return true
}

// Read session data from file: accessToken, refreshToken, publicKey, privateKey
// Note: This is for demo purposes only. To make a secure session, save the credentials in db or some other secure source
func FileReadSession(userId string)(sessionStruct FileSessionStruct, status bool){

	content, err := ioutil.ReadFile("./tmp/" + userId + ".txt")

     if err != nil {
        //   log.Fatal(err)
		  return FileSessionStruct{}, false
     }

    e := json.Unmarshal(content, &sessionStruct)
	if e != nil {
		// log.Println(e.Error())
		return FileSessionStruct{}, false
	}

	return sessionStruct, true
}