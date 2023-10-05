package discovery

import (
	"io/ioutil"
	"net/http"
	"os"
	"sxt-sdks/helpers"
)

// List available namespaces in the blockchain
func ListNamespaces()(namespaces string, errMsg string, status bool){
	apiEndPoint, _ := helpers.ReadEndPointDiscovery()
	tokenEndPoint := apiEndPoint + "/discover/namespace"

	return executeRequest(tokenEndPoint)
}

/* List tables in a given namespace
 Possible scope values -  ALL = all resources, PUBLIC = non-permissioned tables, PRIVATE = tables created by the requesting user*/
func ListTables(namespace, scope string)(tables string, errMsg string, status bool){
	apiEndPoint, _ := helpers.ReadEndPointDiscovery()
	status = true

	re, r := helpers.CheckUpperCase(namespace)
	if !r {
		return "", re, r
	}
	
	tokenEndPoint := apiEndPoint + "/discover/table?scope=" + scope
	if namespace != ""{
		tokenEndPoint += "&namespace=" + namespace
	}

	return executeRequest(tokenEndPoint)
}

// List columns in a given namespace and a table
func ListColumns(namespace, table string)(columns string, errMsg string, status bool){
	apiEndPoint, _ := helpers.ReadEndPointDiscovery()
	status = true

	re, r := helpers.CheckUpperCase(namespace)
	if !r {
		return "", re, r
	}

	re, r = helpers.CheckUpperCase(table)
	if !r {
		return "", re, r
	}

	tokenEndPoint := apiEndPoint + "/discover/table/column?table=" + table + "&namespace=" + namespace

	return executeRequest(tokenEndPoint)
}

// List table index in a given namespace and a table
func ListTableIndex(namespace, table string)(indexes string, errMsg string, status bool){
	apiEndPoint, _ := helpers.ReadEndPointDiscovery()

	status = true

	re, r := helpers.CheckUpperCase(namespace)
	if !r {
		return "", re, r
	}

	re, r = helpers.CheckUpperCase(table)
	if !r {
		return "", re, r
	}

	tokenEndPoint := apiEndPoint + "/discover/table/index?table=" + table + "&namespace=" + namespace

	return executeRequest(tokenEndPoint)
}

// List table primary keys in a given namespace and a table
func ListTablePrimaryKey(namespace, table string)(primaryKeys string, errMsg string, status bool){
	apiEndPoint, _ := helpers.ReadEndPointDiscovery()

	status = true

	re, r := helpers.CheckUpperCase(namespace)
	if !r {
		return "", re, r
	}

	re, r = helpers.CheckUpperCase(table)
	if !r {
		return "", re, r
	}

	tokenEndPoint := apiEndPoint + "/discover/table/primarykey?table=" + table + "&namespace=" + namespace

	return executeRequest(tokenEndPoint)
}

// List table relationships in a given namespace and a table
// Scope can be PRIVATE, PUBLIC, ALL
func ListTableRelations(namespace, scope string)(relations string, errMsg string, status bool){
	apiEndPoint, _ := helpers.ReadEndPointDiscovery()

	status = true

	re, r := helpers.CheckUpperCase(namespace)
	if !r {
		return "", re, r
	}

	tokenEndPoint := apiEndPoint + "/discover/table/relations?scope=" + scope + "&namespace=" + namespace
	return executeRequest(tokenEndPoint)
}

// List primary key references in a given namespace and a table and a column
func ListPrimaryKeyReferences(namespace, table, column string)(primaryKeyReferences string, errMsg string, status bool){
	apiEndPoint, _ := helpers.ReadEndPointDiscovery()
	
	status = true

	re, r := helpers.CheckUpperCase(namespace)
	if !r {
		return "", re, r
	}

	re, r = helpers.CheckUpperCase(table)
	if !r {
		return "", re, r
	}

	re, r = helpers.CheckUpperCase(column)
	if !r {
		return "", re, r
	}

	tokenEndPoint := apiEndPoint + "/discover/refs/primarykey?table=" + table + "&namespace=" + namespace + "&column=" + column
	return executeRequest(tokenEndPoint)
}

// List foreign key references in a given namespace and a table and a column
func ListForeignKeyReferences(namespace, table, column string)(foreignKeyReferences string, errMsg string, status bool){
	apiEndPoint, _ := helpers.ReadEndPointDiscovery()

	status = true

	re, r := helpers.CheckUpperCase(table)
	if !r {
		return "", re, r
	}

	re, r = helpers.CheckUpperCase(column)
	if !r {
		return "", re, r
	}

	tokenEndPoint := apiEndPoint + "/discover/refs/foreignkey?table=" + table + "&namespace=" + namespace + "&column=" + column
	return executeRequest(tokenEndPoint)
}

// List views
// owned values can be a "", 'true', 'false'. All string not boolean
// Both parameters are optional
func ListViews(name, owned string)(foreignKeyReferences string, errMsg string, status bool){
	apiEndPoint, _ := helpers.ReadEndPointDiscovery()
	tokenEndPoint := apiEndPoint + "/discover/views?"

	entryExist := 0

	if name != ""{
		tokenEndPoint += "name=" + name
		entryExist = 1
	}

	if owned != ""{
		if entryExist == 1 {
			tokenEndPoint += "&"
		}
		tokenEndPoint += "owned=" + owned
	}

	return executeRequest(tokenEndPoint)
}

// Master function
func executeRequest(endpoint string)(output string, errMsg string, status bool){
	client := http.Client{}
	req , err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return "", err.Error(), false
	}

	req.Header.Add("Authorization", "Bearer " + os.Getenv("accessToken"))

	res, err := client.Do(req)
	if err != nil {
		return "", err.Error(), false
	}


	defer res.Body.Close()
    body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err.Error(), false
	}

	return string(body), "", true
}