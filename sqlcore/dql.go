package sqlcore

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"sxt-sdks/helpers"
)

// Run all DQL queries
// rowCount is optional
func DQL(sqlText, originApp string, biscuitArray, resources []string, rowCount int) (data []byte, errMsg string, status bool){
	apiEndPoint, _ := helpers.ReadEndPointGeneral()
	tokenEndPoint := apiEndPoint + "/sql/dql"

	
	client := http.Client{}
	var postBody []byte;

	if rowCount > 0 {
		postBody, _ = json.Marshal(map[string]interface{}{
			"biscuits": biscuitArray,
			"resources": resources,
			"sqlText":    sqlText,
			"rowCount": rowCount,
		})
	} else {
		postBody, _ = json.Marshal(map[string]interface{}{
			"biscuits": biscuitArray,
			"resources": resources,
			"sqlText":    sqlText,
		})
	}
	

	responseBody := bytes.NewBuffer(postBody)

	req, err := http.NewRequest("POST", tokenEndPoint, responseBody)
	if err != nil {
		return data, err.Error(), false
	}

	req.Header.Add("Authorization", "Bearer "+os.Getenv("accessToken"))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return data, err.Error(), false
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return data, err.Error(), false
	}

	if res.StatusCode != 200 {
		return data, string(body), false
	}


	return body, "", true
}