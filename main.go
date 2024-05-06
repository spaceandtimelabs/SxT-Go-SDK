package main

import (
	"crypto/ed25519"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/spaceandtimelabs/SxT-Go-SDK/authentication"
	"github.com/spaceandtimelabs/SxT-Go-SDK/helpers"
	"github.com/spaceandtimelabs/SxT-Go-SDK/storage"
	"github.com/spaceandtimelabs/SxT-Go-SDK/utils"
)

// Check the command line arguments
func isFlagPassed(name string) int {
	count := 0
	flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			count = 1
		}
	})
	return count
}



// Main function
func main() {

	fmt.Println("")
	fmt.Println("For exisiting users")
	fmt.Println("Usage: go run main.go -userid=<USERID> -pubkey=<BASE64 STD ENCODED PUBLIC KEY> -privkey=<BASE64 STD ENCODED PRIVATE KEY>")
	fmt.Println("")

	var privateKey ed25519.PrivateKey
	var publicKey ed25519.PublicKey
	var accessToken string


	inputUserid := flag.String("userid", "", "(Optional) SxT userid. But if provided, the remaining values are required")
	inputPubKey := flag.String("pubkey", "", "(Optional) Standard base64 encoded public key. But if provided, the remaining values are required")
	inputPrivKey := flag.String("privkey", "", "(Optional) Standard base64 encoded private key. But if provided, the remaining values are required")
	flag.Parse()


	/*************************************
	// Authentication Block
	*************************************/

	/* AUTH BLOCK STARTS */
	totalFlags := isFlagPassed("userid") + isFlagPassed("pubkey") + isFlagPassed("privkey")

	if totalFlags < 3 && totalFlags > 0 {
		fmt.Println("=== Missing input values. Stopping program ===")
		return
	}

	if isFlagPassed("userid")+isFlagPassed("pubkey")+isFlagPassed("privkey") == 3 {

		if len(*inputUserid) > 0 || len(*inputPubKey) > 0 || len(*inputPrivKey) > 0 {
			accessToken, _, privateKey, publicKey, _ = utils.Authenticate(*inputUserid, *inputPubKey, *inputPrivKey)
			fmt.Println("=== Existing Login from user input ===", accessToken)

		} else {
			fmt.Println("=== Empty input values. Stopping program ===")
			return
		}
	} else {
		userId, _ := helpers.ReadUserId()
		sessionData, status := storage.FileReadSession(userId)

		if !status {

			emptyString := ""
			accessToken, _, privateKey, publicKey, _ = utils.Authenticate(emptyString, emptyString, emptyString)
			fmt.Println("=== New Login. Creating new session ===")

		} else {

			validateTokenStatus := authentication.ValidateToken(sessionData.AccessToken)

			if !validateTokenStatus {

				refreshTokenStruct, refreshTokenStatus := authentication.RefreshToken(sessionData.RefreshToken)

				if refreshTokenStatus {
					fmt.Println("=== Invalid session on session.txt file. Using refresh token ===")

					privateKey = sessionData.PrivateKey
					publicKey = sessionData.PublicKey
					accessToken = refreshTokenStruct.AccessToken

					writeStatus := storage.FileWriteSession(userId, refreshTokenStruct.AccessToken, refreshTokenStruct.RefreshToken, sessionData.PrivateKey, sessionData.PublicKey)
					if !writeStatus {
						log.Fatal("Invalid login. Change login credentials")
					}

				} else {
					fmt.Println("=== Invalid session on session.txt file. Issuing new token ===")
					emptyString := ""
					accessToken, _, privateKey, publicKey, _ = utils.Authenticate(emptyString, emptyString, emptyString)
				}
			} else {
				fmt.Println("=== Login using session.txt file ===")
				privateKey = sessionData.PrivateKey
				publicKey = sessionData.PublicKey
				accessToken = sessionData.AccessToken
			}
		}
	}

	// Important : Set the accessToken in Environment variable to access in later
	os.Setenv("accessToken", accessToken)

	/* AUTH BLOCK ENDS */


	// Call SQL functions
	utils.SQLAPIs(privateKey, publicKey)
	utils.DiscoveryAPIs()
	
}
