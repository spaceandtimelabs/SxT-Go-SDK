package storage

import (
	"context"
	"crypto/ed25519"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

type AwsSessionStruct struct{
	AccessToken string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	PrivateKey ed25519.PrivateKey `json:"privateKey"`
	PublicKey ed25519.PublicKey `json:"publicKey"`
}

// Write session data to aws secrets manager: accessToken, refreshToken, publicKey, privateKey
// Filename to be stored as `userId` (provided by SxT)
func AwsWriteSession(userId, accessToken, refreshToken string, privateKey ed25519.PrivateKey, publicKey ed25519.PublicKey)( status bool){

	if accessToken == "" || refreshToken == "" {
		return false
	}
	
 	sessionStruct := AwsSessionStruct{
		AccessToken: accessToken,
		RefreshToken: refreshToken,
		PublicKey: publicKey,
		PrivateKey: privateKey,
	}

	sessionData, e := json.Marshal(sessionStruct)

	if e != nil {
		return false
	} else {
		
		sess, _ := awsAuthenticate()
		
		input := &secretsmanager.CreateSecretInput{
			Description:        aws.String("Credentials for "+ userId),
			Name:               aws.String(userId), // In live, you can put the sxt username (userId) here
			SecretString:       aws.String(string(sessionData)),
		}

		svc := secretsmanager.NewFromConfig(sess)
		result, err := svc.CreateSecret(context.TODO(), input)

		if err != nil {
			return false
		}

		fmt.Println("WRITE Session", result)
	}

	return true
}

// Read session data from aws secrets manager: accessToken, refreshToken, publicKey, privateKey
// Filename to be retrieved as `userId` (provided by SxT)
func AwsReadSession(userId string)(sessionStruct AwsSessionStruct, status bool){

	sess, _ := awsAuthenticate()

	input := &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(userId), // In live, you can put the sxt username (userId) here
	}

	svc := secretsmanager.NewFromConfig(sess)
	result, err := svc.GetSecretValue(context.TODO(), input)
	if err != nil {
		return AwsSessionStruct{}, false
	}

	// fmt.Println("READ", *result.SecretString)
    e := json.Unmarshal([]byte(*result.SecretString), &sessionStruct)
	if e != nil {
		return AwsSessionStruct{}, false
	}

	return sessionStruct, true
}


// Update session data to aws secrets manager: accessToken, refreshToken, publicKey, privateKey
// Filename to be updated as `userId` (provided by SxT)
func AwsUpdateSession(userId, accessToken, refreshToken string, privateKey ed25519.PrivateKey, publicKey ed25519.PublicKey)( status bool){

	if accessToken == "" || refreshToken == "" {
		return false
	}
	
 	sessionStruct := AwsSessionStruct{
		AccessToken: accessToken,
		RefreshToken: refreshToken,
		PublicKey: publicKey,
		PrivateKey: privateKey,
	}

	sessionData, e := json.Marshal(sessionStruct)

	if e != nil {
		return false
	} else {
		
		sess, _ := awsAuthenticate()
		
		input := &secretsmanager.PutSecretValueInput{
			SecretId:           aws.String(userId), // In live, you can put the sxt username (userId) here
			SecretString:       aws.String(string(sessionData)),
		}

		svc := secretsmanager.NewFromConfig(sess)
		result, err := svc.PutSecretValue(context.TODO(), input)

		if err != nil {
			return false
		}

		fmt.Println("UPDATE Session", result)
	}

	return true
}

// aws authenticate
func awsAuthenticate()(cfg aws.Config, err error){

	cfg , err = config.LoadDefaultConfig(context.TODO());
	return cfg, err
}