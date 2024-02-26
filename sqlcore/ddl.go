package sqlcore

import (
	"crypto/ed25519"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Create a new table on a given namespace.
// accessType: can be public, permissioned or encrypted. Read more here https://docs.spaceandtime.io/docs/secure-your-table
func CreateTable(sqlText, accessType, originApp string, biscuitArray []string, publicKey ed25519.PublicKey) (errMsg string, status bool) {
	sqlTextWithConfiguration := fmt.Sprintf("%s WITH \"public_key=%x,access_type=%s\"", sqlText, publicKey, accessType)

	return DDL(sqlTextWithConfiguration, originApp, biscuitArray)
}

// DDL queries for ALTER and DROP
func DDL(sqlText, originApp string, biscuitArray []string) (errMsg string, status bool) {
	request, err := createResourceConfigurationRequest(sqlText, originApp, biscuitArray)
	if err != nil {
		return err.Error(), false
	}

	client := http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return err.Error(), false
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err.Error(), false
	}

	if response.StatusCode != 200 {
		return string(body), false
	}

	return "", true

}

func createResourceConfigurationRequest(sqlText, originApp string, biscuitArray []string) (request *http.Request, err error) {
	postBody, _ := json.Marshal(map[string]interface{}{
		"biscuits": biscuitArray,
		"sqlText":  sqlText,
	})

	return createRequest("ddl", originApp, postBody)
}

// Create a new schema
func CreateSchema(sqlText, originApp string, biscuitArray []string) (errMsg string, status bool) {
	return DDL(sqlText, originApp, biscuitArray)
}
