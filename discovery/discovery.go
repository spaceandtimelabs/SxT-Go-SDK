package discovery

import (
	"io/ioutil"
	"net/http"
	"os"
	"sxt-sdks/helpers"
)

// List available namespaces in the blockchain
func ListSchemas(scope, searchPattern string)(schemas string, errMsg string, status bool){
	apiEndPoint, _ := helpers.ReadEndPointDiscovery()
	tokenEndPoint := apiEndPoint + "/discover/schema?scope=" + scope
	if searchPattern != ""{
		tokenEndPoint += tokenEndPoint + "&searchPattern=" + searchPattern
	}

	return executeRequest(tokenEndPoint)
}

/* List tables in a given schema
 Possible scope values -  ALL = all resources, PUBLIC = non-permissioned tables, PRIVATE = tables created by the requesting user*/
func ListTables(schema, scope, searchPattern string)(tables string, errMsg string, status bool){
	apiEndPoint, _ := helpers.ReadEndPointDiscovery()
	status = true

	re, r := helpers.CheckUpperCase(schema)
	if !r {
		return "", re, r
	}
	
	tokenEndPoint := apiEndPoint + "/discover/table?scope=" + scope
	if schema != ""{
		tokenEndPoint += "&schema=" + schema
	}

	if searchPattern != ""{
		tokenEndPoint += tokenEndPoint + "&searchPattern=" + searchPattern
	}

	return executeRequest(tokenEndPoint)
}

// List columns in a given schema and a table
func ListColumns(schema, table string)(columns string, errMsg string, status bool){
	apiEndPoint, _ := helpers.ReadEndPointDiscovery()
	status = true

	re, r := helpers.CheckUpperCase(schema)
	if !r {
		return "", re, r
	}

	re, r = helpers.CheckUpperCase(table)
	if !r {
		return "", re, r
	}

	tokenEndPoint := apiEndPoint + "/discover/table/column?table=" + table + "&schema=" + schema

	return executeRequest(tokenEndPoint)
}

// List table index in a given schema and a table
func ListTableIndex(schema, table string)(indexes string, errMsg string, status bool){
	apiEndPoint, _ := helpers.ReadEndPointDiscovery()

	status = true

	re, r := helpers.CheckUpperCase(schema)
	if !r {
		return "", re, r
	}

	re, r = helpers.CheckUpperCase(table)
	if !r {
		return "", re, r
	}

	tokenEndPoint := apiEndPoint + "/discover/table/index?table=" + table + "&schema=" + schema

	return executeRequest(tokenEndPoint)
}

// List table primary keys in a given schema and a table
func ListTablePrimaryKey(schema, table string)(primaryKeys string, errMsg string, status bool){
	apiEndPoint, _ := helpers.ReadEndPointDiscovery()

	status = true

	re, r := helpers.CheckUpperCase(schema)
	if !r {
		return "", re, r
	}

	re, r = helpers.CheckUpperCase(table)
	if !r {
		return "", re, r
	}

	tokenEndPoint := apiEndPoint + "/discover/table/primarykey?table=" + table + "&schema=" + schema

	return executeRequest(tokenEndPoint)
}

// List table relationships in a given schema and a table
// Scope can be PRIVATE, PUBLIC, ALL
func ListTableRelations(schema, scope string)(relations string, errMsg string, status bool){
	apiEndPoint, _ := helpers.ReadEndPointDiscovery()

	status = true

	re, r := helpers.CheckUpperCase(schema)
	if !r {
		return "", re, r
	}

	tokenEndPoint := apiEndPoint + "/discover/table/relations?scope=" + scope + "&schema=" + schema
	return executeRequest(tokenEndPoint)
}

// List primary key references in a given schema and a table and a column
func ListPrimaryKeyReferences(schema, table, column string)(primaryKeyReferences string, errMsg string, status bool){
	apiEndPoint, _ := helpers.ReadEndPointDiscovery()
	
	status = true

	re, r := helpers.CheckUpperCase(schema)
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

	tokenEndPoint := apiEndPoint + "/discover/refs/primarykey?table=" + table + "&schema=" + schema + "&column=" + column
	return executeRequest(tokenEndPoint)
}

// List foreign key references in a given schema and a table and a column
func ListForeignKeyReferences(schema, table, column string)(foreignKeyReferences string, errMsg string, status bool){
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

	tokenEndPoint := apiEndPoint + "/discover/refs/foreignkey?table=" + table + "&schema=" + schema + "&column=" + column
	return executeRequest(tokenEndPoint)
}

// List Blockchains
func ListBlockchains( )(blockchains string, errMsg string, status bool){
	apiEndPoint, _ := helpers.ReadEndPointDiscovery()
	tokenEndPoint := apiEndPoint + "/discover/blockchains"

	return executeRequest(tokenEndPoint)
}


// List Blockchains
func ListBlockchainSchemas(chainId string)(blockchainSchema string, errMsg string, status bool){
	apiEndPoint, _ := helpers.ReadEndPointDiscovery()
	
	tokenEndPoint := apiEndPoint + "/discover/blockchains/" + chainId + "/schemas"

	return executeRequest(tokenEndPoint)
}

// List Blockchain Information
func ListBlockchainInformation(chainId string)(blockchainInformation string, errMsg string, status bool){
	apiEndPoint, _ := helpers.ReadEndPointDiscovery()
	
	tokenEndPoint := apiEndPoint + "/discover/blockchains/" + chainId + "/meta"

	return executeRequest(tokenEndPoint)
}

// List views
// owned values can be a "", 'true', 'false'. All string not boolean
// Both parameters are optional
func ListViews(name, owned string)(views string, errMsg string, status bool){
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