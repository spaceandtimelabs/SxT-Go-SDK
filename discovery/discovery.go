package discovery

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"sxt-sdks/helpers"
)

// List available namespaces in the blockchain
func ListSchemas(scope, searchPattern string) (schemas string, errMsg string, status bool) {
	tokenEndPoint := helpers.GetDiscoverEndpoint("schema") + "?scope=" + scope

	if searchPattern != "" {
		tokenEndPoint += "&searchPattern=" + searchPattern
	}

	return executeRequest(tokenEndPoint)
}

/*
List tables in a given schema
Possible scope values -  ALL = all resources, PUBLIC = non-permissioned tables, PRIVATE = tables created by the requesting user
*/
func ListTables(schema, scope, searchPattern string) (tables string, errMsg string, status bool) {
	re, r := helpers.CheckUpperCase(schema)
	if !r {
		return "", re, r
	}

	tableEndpoint := helpers.GetDiscoverEndpoint("table")
	tokenEndPoint := fmt.Sprintf("%s?scope=%s", tableEndpoint, scope)
	if schema != "" {
		tokenEndPoint += "&schema=" + schema
	}

	if searchPattern != "" {
		tokenEndPoint += "&searchPattern=" + searchPattern
	}

	return executeRequest(tokenEndPoint)
}

// List columns in a given schema and a table
func ListColumns(schema, table string) (columns string, errMsg string, status bool) {
	return listTableInfo("column", schema, table)
}

// List table index in a given schema and a table
func ListTableIndex(schema, table string) (indexes string, errMsg string, status bool) {
	return listTableInfo("index", schema, table)
}

// List table primary keys in a given schema and a table
func ListTablePrimaryKey(schema, table string) (primaryKeys string, errMsg string, status bool) {
	return listTableInfo("primarykey", schema, table)
}

// List table primary keys in a given schema and a table
func listTableInfo(infoType, schema, table string) (info string, errMsg string, status bool) {
	queryParameters := []string{schema, table}

	for _, field := range queryParameters {
		message, isUpperCase := helpers.CheckUpperCase(field)
		if !isUpperCase {
			return "", message, isUpperCase
		}
	}

	tableEndpoint := helpers.GetDiscoverEndpoint("table")
	tokenEndPoint := fmt.Sprintf("%s/%s?schema=%s&table=%s", tableEndpoint, infoType, schema, table)

	return executeRequest(tokenEndPoint)

}

// List table relationships in a given schema and a table
// Scope can be PRIVATE, PUBLIC, ALL
func ListTableRelations(schema, scope string) (relations string, errMsg string, status bool) {
	re, r := helpers.CheckUpperCase(schema)
	if !r {
		return "", re, r
	}

	tableEndpoint := helpers.GetDiscoverEndpoint("table")
	tokenEndPoint := fmt.Sprintf("%s/relations?schema=%s&scope=%s", tableEndpoint, schema, scope)

	return executeRequest(tokenEndPoint)
}

// List primary key references in a given schema and a table and a column
func ListPrimaryKeyReferences(schema, table, column string) (primaryKeyReferences string, errMsg string, status bool) {
	return listKeyReferences("primary", schema, table, column)
}

// List foreign key references in a given schema and a table and a column
func ListForeignKeyReferences(schema, table, column string) (foreignKeyReferences string, errMsg string, status bool) {
	return listKeyReferences("foreign", schema, table, column)
}

func listKeyReferences(keyReferenceType, schema, table, column string) (keyReferences string, errMsg string, status bool) {
	queryParameters := []string{schema, table, column}

	for _, field := range queryParameters {
		message, isUpperCase := helpers.CheckUpperCase(field)
		if !isUpperCase {
			return "", message, isUpperCase
		}
	}

	referenceKeyEndpoint := helpers.GetDiscoverEndpoint("refs")
	tokenEndPoint := fmt.Sprintf("%s/%skey?schema=%s&table=%s&column=%s", referenceKeyEndpoint, keyReferenceType, schema, table, column)
	return executeRequest(tokenEndPoint)
}

// List Blockchains
func ListBlockchains() (blockchains string, errMsg string, status bool) {
	return listBlockchainInfo("", "")
}

// List Blockchains
func ListBlockchainSchemas(chainId string) (blockchainSchema string, errMsg string, status bool) {
	return listBlockchainInfo(chainId, "schemas")
}

// List Blockchain Information
func ListBlockchainInformation(chainId string) (blockchainInformation string, errMsg string, status bool) {
	return listBlockchainInfo(chainId, "meta")
}

func listBlockchainInfo(chainId, infoType string) (blockchainInformation string, errMsg string, status bool) {
	discoverBlockchainsEndpoint := helpers.GetDiscoverEndpoint("blockchains")

	if chainId == "" {
		return executeRequest(discoverBlockchainsEndpoint)
	}

	segments := []string{discoverBlockchainsEndpoint, chainId, infoType}
	tokenEndPoint := strings.Join(segments, "/")

	return executeRequest(tokenEndPoint)
}

// List views
// owned values can be a "", 'true', 'false'. All string not boolean
// Both parameters are optional
func ListViews(name, owned string) (views string, errMsg string, status bool) {
	tokenEndPoint := helpers.GetDiscoverEndpoint("views") + "?"
	entryExists := false

	if name != "" {
		tokenEndPoint += "name=" + name
		entryExists = true
	}

	if owned != "" {
		if entryExists {
			tokenEndPoint += "&"
		}
		tokenEndPoint += "owned=" + owned
	}

	return executeRequest(tokenEndPoint)
}

// Master function
func executeRequest(endpoint string) (output string, errMsg string, status bool) {
	client := http.Client{}
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return "", err.Error(), false
	}

	req.Header.Add("Authorization", "Bearer "+os.Getenv("accessToken"))

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
