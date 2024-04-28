package sqlcore

import (
	"encoding/json"
	"io"
	"net/http"
)

// Run all DQL queries
// rowCount is optional
func DQL(sqlText, originApp string, biscuitArray, resources []string, rowCount int) (data []byte, errMsg string, status bool) {
	request, err := createQueryExecutionRequest(sqlText, originApp, biscuitArray, resources)
	if err != nil {
		return data, err.Error(), false
	}

	client := http.Client{}
	res, err := client.Do(request)
	if err != nil {
		return data, err.Error(), false
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return data, err.Error(), false
	}

	if res.StatusCode != 200 {
		return data, string(body), false
	}

	return body, "", true
}

func createQueryExecutionRequest(sqlText, originApp string, biscuitArray, resources []string) (request *http.Request, err error) {
	postBody, _ := json.Marshal(map[string]interface{}{
		"biscuits":  biscuitArray,
		"resources": resources,
		"sqlText":   sqlText,
	})

	return createRequest("dql", originApp, postBody)
}
