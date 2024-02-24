package sqlcore

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// Run all DML queries
func DML(sqlText, originApp string, biscuitArray []string, resources []string) (errMsg string, status bool) {
	request, err := createDataModificationRequest(sqlText, originApp, biscuitArray, resources)
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

func createDataModificationRequest(sqlText, originApp string, biscuitArray, resources []string) (request *http.Request, err error) {
	postBody, _ := json.Marshal(map[string]interface{}{
		"biscuits":  biscuitArray,
		"resources": resources,
		"sqlText":   sqlText,
	})

	return createRequest("dml", originApp, postBody)
}
