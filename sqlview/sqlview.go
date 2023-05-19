package sqlview

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"sxt-sdks/helpers"
)
type ParametersRequest struct{
	Name string `json:"Name"`
	Type string `json:"Type"`
}


// Create sql view
// parameterRequest is optional
func Create(resourceId, viewName, viewText, description string, publish bool, parametersRequest []ParametersRequest) (output, errMsg string, status bool){
	apiEndPoint, _ := helpers.ReadEndPoint()
	tokenEndPoint := apiEndPoint + "/sql/views"

	re, r := helpers.CheckUpperCaseResource(resourceId)
	if !r {
		return "", re, r
	}

	var postBody []byte
	if len(parametersRequest) > 0{
		pr, _ := json.Marshal(parametersRequest)

		postBody, _ = json.Marshal(map[string]interface{}{
			"viewName": viewName,
			"viewText":    viewText,
			"resourceId": resourceId,
			"description": description,
			"publish": publish,
			"parametersRequest": string(pr),
		})
	} else {
		postBody, _ = json.Marshal(map[string]interface{}{
			"viewName": viewName,
			"viewText":    viewText,
			"resourceId": resourceId,
			"description": description,
			"publish": publish,
		})
	}
	
	client := http.Client{}
	responseBody := bytes.NewBuffer(postBody)

	req, err := http.NewRequest("POST", tokenEndPoint, responseBody)
	if err != nil {
		return "", err.Error(), false
	}

	req.Header.Add("Authorization", "Bearer "+os.Getenv("accessToken"))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return "", err.Error(), false
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err.Error(), false
	}

	if res.StatusCode != 201 {
		return "", string(body), false
	}

	return string(body), "", true
}

// Execute a view
// parameterRequest is optional
func Execute(viewName string, parametersRequest []ParametersRequest) (errMsg string, status bool){
	apiEndPoint, _ := helpers.ReadEndPoint()
	tokenEndPoint := apiEndPoint + "/sql/views/"+ url.QueryEscape(viewName)

	paramString := ""

	if len(parametersRequest) > 0 {
		for _, param := range parametersRequest {
			paramString += param.Name + "=" + param.Type + "&"
		}
		paramString = paramString[:len(paramString)-1]

		tokenEndPoint += "?params=" + paramString
	}


	client := http.Client{}
	req , err := http.NewRequest("GET", tokenEndPoint, nil)
	if err != nil {
		return err.Error(), false
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", "Bearer " + os.Getenv("accessToken"))

	res, err := client.Do(req)
	if err != nil {
		return err.Error(), false
	}


	defer res.Body.Close()
    _, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return err.Error(), false
	}

	return "", true
}

// Update sql view
// parameterRequest is optional
func Update(resourceId, viewName, viewText, description string, publish bool, parametersRequest []ParametersRequest) (output, errMsg string, status bool){
	apiEndPoint, _ := helpers.ReadEndPoint()
	tokenEndPoint := apiEndPoint + "/sql/views/"+ viewName

	re, r := helpers.CheckUpperCaseResource(resourceId)
	if !r {
		return "", re, r
	}

	var postBody []byte
	if len(parametersRequest) > 0{
		pr, _ := json.Marshal(parametersRequest)

		postBody, _ = json.Marshal(map[string]interface{}{
			"viewName": viewName,
			"viewText":    viewText,
			"resourceId": resourceId,
			"description": description,
			"publish": publish,
			"parametersRequest": string(pr),
		})
	} else {
		postBody, _ = json.Marshal(map[string]interface{}{
			"viewName": viewName,
			"viewText":    viewText,
			"resourceId": resourceId,
			"description": description,
			"publish": publish,
		})
	}
	
	client := http.Client{}
	responseBody := bytes.NewBuffer(postBody)

	req, err := http.NewRequest("PUT", tokenEndPoint, responseBody)
	if err != nil {
		return "", err.Error(), false
	}

	req.Header.Add("Authorization", "Bearer "+os.Getenv("accessToken"))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return "", err.Error(), false
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err.Error(), false
	}

	if res.StatusCode != 204 {
		return "", string(body), false
	}

	return string(body), "", true
}


// Delete sql view
// parameterRequest is optional
func Delete(viewName string) (output, errMsg string, status bool){
	apiEndPoint, _ := helpers.ReadEndPoint()
	tokenEndPoint := apiEndPoint + "/sql/views/"+ viewName

	client := http.Client{}

	req, err := http.NewRequest("DELETE", tokenEndPoint, nil)
	if err != nil {
		return "", err.Error(), false
	}

	req.Header.Add("Authorization", "Bearer "+os.Getenv("accessToken"))
	req.Header.Add("Accept", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return "", err.Error(), false
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err.Error(), false
	}

	if res.StatusCode != 204 {
		return "", string(body), false
	}

	return string(body), "", true
}