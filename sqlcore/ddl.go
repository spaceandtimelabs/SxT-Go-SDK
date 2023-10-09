package sqlcore

import (
	"bytes"
	"crypto/ed25519"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sxt-sdks/helpers"
)

// Create a new table on a given namespace.
// accessType: can be public, permissioned or encrypted. Read more here https://docs.spaceandtime.io/docs/secure-your-table
func CreateTable(sqlText, accessType, originApp string, biscuitArray []string, publicKey ed25519.PublicKey)(errMsg string, status bool) {
	apiEndPoint, _ := helpers.ReadEndPointGeneral()
	tokenEndPoint := apiEndPoint + "/sql/ddl"

	client := http.Client{}
	postBody, _ := json.Marshal(map[string]interface{}{
		"biscuits": biscuitArray,
		"sqlText": sqlText + " WITH \"public_key=" + fmt.Sprintf("%x",publicKey) + ",access_type=" + accessType + "\"",
	})

	responseBody := bytes.NewBuffer(postBody)

	req, err := http.NewRequest("POST", tokenEndPoint, responseBody)
	if err != nil {
		return err.Error(), false
	}

	req.Header.Add("Authorization", "Bearer "+os.Getenv("accessToken"))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("originApp", originApp)

	res, err := client.Do(req)
	if err != nil {
		return err.Error(), false
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err.Error(), false
	}

	if res.StatusCode != 200 {
		return string(body), false
	} 

	return "", true
}

// DDL queries for ALTER and DROP
func DDL(sqlText, originApp  string, biscuitArray []string) (errMsg string, status bool){
	apiEndPoint, _ := helpers.ReadEndPointGeneral()
	tokenEndPoint := apiEndPoint + "/sql/ddl"

	client := http.Client{}
	postBody, _ := json.Marshal(map[string]interface{}{
		"biscuits": biscuitArray,
		"sqlText":    sqlText,
	})

	responseBody := bytes.NewBuffer(postBody)

	req, err := http.NewRequest("POST", tokenEndPoint, responseBody)
	if err != nil {
		return err.Error(), false
	}

	req.Header.Add("Authorization", "Bearer "+os.Getenv("accessToken"))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("originApp", originApp)

	res, err := client.Do(req)
	if err != nil {
		return err.Error(), false
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err.Error(), false
	}

	if res.StatusCode != 200 {
		return string(body), false
	}

	return "", true
	
}

// Create a new schema
func CreateSchema(sqlText, originApp  string, biscuitArray []string)(errMsg string, status bool) {
	apiEndPoint, _ := helpers.ReadEndPointGeneral()
	tokenEndPoint := apiEndPoint + "/sql/ddl"

	client := http.Client{}
	postBody, _ := json.Marshal(map[string]interface{}{
		"biscuits": biscuitArray,
		"sqlText": sqlText,
	})

	responseBody := bytes.NewBuffer(postBody)

	req, err := http.NewRequest("POST", tokenEndPoint, responseBody)
	if err != nil {
		return err.Error(), false
	}

	req.Header.Add("Authorization", "Bearer "+os.Getenv("accessToken"))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("originApp", originApp)

	res, err := client.Do(req)
	if err != nil {
		return err.Error(), false
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err.Error(), false
	}

	if res.StatusCode != 200 {
		return string(body), false
	} 

	return "", true
}