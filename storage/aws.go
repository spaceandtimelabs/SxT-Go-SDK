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

type AwsSessionStruct struct {
	AccessToken  string             `json:"accessToken"`
	RefreshToken string             `json:"refreshToken"`
	PrivateKey   ed25519.PrivateKey `json:"privateKey"`
	PublicKey    ed25519.PublicKey  `json:"publicKey"`
}

// Write session data to aws secrets manager: accessToken, refreshToken, publicKey, privateKey
// Filename to be stored as `userId` (provided by SxT)
func AwsWriteSession(userId, accessToken, refreshToken string, privateKey ed25519.PrivateKey, publicKey ed25519.PublicKey) (status bool) {
	if accessToken == "" || refreshToken == "" {
		return false
	}

	sessionStruct := AwsSessionStruct{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		PublicKey:    publicKey,
		PrivateKey:   privateKey,
	}

	sessionData, err := json.Marshal(sessionStruct)

	if err != nil {
		return false
	} else {
		secretsManagerClient := getSecretsManagerClient()
		input := &secretsmanager.CreateSecretInput{
			Description:  aws.String("Credentials for " + userId),
			Name:         aws.String(userId), // In live, you can put the sxt username (userId) here
			SecretString: aws.String(string(sessionData)),
		}

		secret, err := secretsManagerClient.CreateSecret(context.TODO(), input)
		if err != nil {
			return false
		}

		fmt.Println("WRITE Session", secret)
	}

	return true
}

// Read session data from aws secrets manager: accessToken, refreshToken, publicKey, privateKey
// Filename to be retrieved as `userId` (provided by SxT)
func AwsReadSession(userId string) (sessionStruct AwsSessionStruct, status bool) {
	secretsManagerClient := getSecretsManagerClient()
	input := &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(userId), // In live, you can put the sxt username (userId) here
	}

	secret, err := secretsManagerClient.GetSecretValue(context.TODO(), input)
	if err != nil || json.Unmarshal([]byte(*secret.SecretString), &sessionStruct) != nil {
		return AwsSessionStruct{}, false
	}

	return sessionStruct, true
}

// Update session data to aws secrets manager: accessToken, refreshToken, publicKey, privateKey
// Filename to be updated as `userId` (provided by SxT)
func AwsUpdateSession(userId, accessToken, refreshToken string, privateKey ed25519.PrivateKey, publicKey ed25519.PublicKey) (status bool) {
	if accessToken == "" || refreshToken == "" {
		return false
	}

	sessionStruct := AwsSessionStruct{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		PublicKey:    publicKey,
		PrivateKey:   privateKey,
	}

	sessionData, err := json.Marshal(sessionStruct)
	if err != nil {
		return false
	} else {
		secretsManagerClient := getSecretsManagerClient()
		input := &secretsmanager.PutSecretValueInput{
			SecretId:     aws.String(userId), // In live, you can put the sxt username (userId) here
			SecretString: aws.String(string(sessionData)),
		}

		secret, err := secretsManagerClient.PutSecretValue(context.TODO(), input)
		if err != nil {
			return false
		}

		fmt.Println("UPDATE Session", secret)
	}

	return true
}

func getSecretsManagerClient() *secretsmanager.Client {
	session, _ := awsAuthenticate()

	return secretsmanager.NewFromConfig(session)
}

func awsAuthenticate() (cfg aws.Config, err error) {
	return config.LoadDefaultConfig(context.TODO())
}
