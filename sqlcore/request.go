package sqlcore

import (
	"bytes"
	"fmt"
	"net/http"
	"os"

	"spaceandtime.io/sxt-sdk/helpers"
)

func createRequest(requestType, originApp string, postBody []byte) (request *http.Request, err error) {
	tokenEndPoint := helpers.GetSqlEndpoint(requestType)
	responseBody := bytes.NewBuffer(postBody)

	request, err = http.NewRequest("POST", tokenEndPoint, responseBody)
	if err != nil {
		return nil, err
	}

	bearerToken := fmt.Sprintf("Bearer %s", os.Getenv("accessToken"))
	contentType := "application/json"
	request.Header.Add("Authorization", bearerToken)
	request.Header.Add("Content-Type", contentType)
	request.Header.Add("Accept", contentType)
	request.Header.Add("originApp", originApp)

	return request, nil
}
