package main

import (
	"crypto/ed25519"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"sxt-sdks/authentication"
	"sxt-sdks/authorization"
	"sxt-sdks/discovery"
	"sxt-sdks/helpers"
	"sxt-sdks/sqlcore"
	"sxt-sdks/sqlview"
	"sxt-sdks/storage"
)

// Check the command line arguments
func isFlagPassed(name string) int {
	count := 0
	flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			count = 1
		}
	})
	return count
}

// New Authentication.
// This method generates new accessToken, refreshToken, privateKey, and publicKey
func authenticate(inputUserId, inputPublicKey, inputPrivateKey string) (accessToken, refreshToken string, privateKey ed25519.PrivateKey, publicKey ed25519.PublicKey) {

	/*************************************
	// AUTHENTICATION APIS
	*************************************/

	userId, _ := helpers.ReadUserId()
	joinCode, _ := helpers.ReadJoinCode()

	var pubkey ed25519.PublicKey
	var privkey ed25519.PrivateKey
	var e error

	var authCodeStruct authentication.AuthCodeStruct
	var tokenStruct authentication.TokenStruct
	var sessionStruct storage.FileSessionStruct
	var sessionStatus bool

	if inputUserId != "" && inputPublicKey != "" && inputPrivateKey != "" {
		pubkey, e = base64.StdEncoding.DecodeString(inputPublicKey)
		if e != nil {
			log.Fatal("base64 std encoded public key expected")
		}

		privkey, e = base64.StdEncoding.DecodeString(inputPrivateKey)
		if e != nil {
			log.Fatal("base64 std encoded private key expected")
		}

		// Key compatibility
		if len(privkey) == len(pubkey) {
			var tmpPvtKey []byte
			for idx := range privkey {
				tmpPvtKey = append(tmpPvtKey, privkey[idx])
			}

			for idx := range pubkey {
				tmpPvtKey = append(tmpPvtKey, pubkey[idx])
			}
			privkey = tmpPvtKey
		}

		userId = inputUserId
	} else {
		sessionStruct, sessionStatus = storage.FileReadSession(userId)

		if !sessionStatus {
			// Generate Private and Public keys
			pubkey, privkey = helpers.CreateKey()
		} else {
			pubkey = sessionStruct.PublicKey
			privkey = sessionStruct.PrivateKey
		}
	}

	// Get auth code
	authCode := authentication.GenerateAuthCode(userId, joinCode)
	json.Unmarshal([]byte(authCode), &authCodeStruct)

	// Get Keys
	encodedSignature, base64PublicKey := authentication.GenerateKeys(authCodeStruct.AuthCode, pubkey, privkey)

	// Get Token
	tokenJson := authentication.GenerateToken(userId, authCodeStruct.AuthCode, encodedSignature, base64PublicKey)
	json.Unmarshal([]byte(tokenJson), &tokenStruct)

	// fmt.Println(userId, tokenStruct.AccessToken, tokenStruct.RefreshToken, privkey, pubkey)

	writeStatus := storage.FileWriteSession(userId, tokenStruct.AccessToken, tokenStruct.RefreshToken, privkey, pubkey)
	if !writeStatus {
		log.Fatal("Invalid login. Change login credentials")
	}

	// Logout
	// authentication.Logout()
	return tokenStruct.AccessToken, tokenStruct.RefreshToken, privkey, publicKey
}

// Main function
func main() {

	fmt.Println("")
	fmt.Println("For exisiting users")
	fmt.Println("Usage: go run main.go -userid=<USERID> -pubkey=<BASE64 STD ENCODED PUBLIC KEY> -privkey=<BASE64 STD ENCODED PRIVATE KEY>")
	fmt.Println("")

	var privateKey ed25519.PrivateKey
	var publicKey ed25519.PublicKey
	var accessToken string

	var sxtBiscuitCapabilities []authorization.SxTBiscuitStruct

	inputUserid := flag.String("userid", "", "(Optional) SxT userid. But if provided, the remaining values are required")
	inputPubKey := flag.String("pubkey", "", "(Optional) Standard base64 encoded public key. But if provided, the remaining values are required")
	inputPrivKey := flag.String("privkey", "", "(Optional) Standard base64 encoded private key. But if provided, the remaining values are required")
	flag.Parse()

	/*************************************
	// Authentication Block
	*************************************/

	/* AUTH BLOCK STARTS */
	totalFlags := isFlagPassed("userid") + isFlagPassed("pubkey") + isFlagPassed("privkey")

	if totalFlags < 3 && totalFlags > 0 {
		fmt.Println("=== Missing input values. Stopping program ===")
		return
	}

	if isFlagPassed("userid")+isFlagPassed("pubkey")+isFlagPassed("privkey") == 3 {

		if len(*inputUserid) > 0 || len(*inputPubKey) > 0 || len(*inputPrivKey) > 0 {
			accessToken, _, privateKey, publicKey = authenticate(*inputUserid, *inputPubKey, *inputPrivKey)
			fmt.Println("=== Existing Login from user input ===")

		} else {
			fmt.Println("=== Empty input values. Stopping program ===")
			return
		}
	} else {
		userId, _ := helpers.ReadUserId()
		sessionData, status := storage.FileReadSession(userId)

		if !status {

			emptyString := ""
			accessToken, _, privateKey, publicKey = authenticate(emptyString, emptyString, emptyString)
			fmt.Println("=== New Login. Creating new session ===")

		} else {

			validateTokenStatus := authentication.ValidateToken(sessionData.AccessToken)

			if !validateTokenStatus {

				refreshTokenStruct, refreshTokenStatus := authentication.RefreshToken(sessionData.RefreshToken)

				if refreshTokenStatus {
					fmt.Println("=== Invalid session on session.txt file. Using refresh token ===")

					privateKey = sessionData.PrivateKey
					publicKey = sessionData.PublicKey
					accessToken = refreshTokenStruct.AccessToken

					writeStatus := storage.FileWriteSession(userId, refreshTokenStruct.AccessToken, refreshTokenStruct.RefreshToken, sessionData.PrivateKey, sessionData.PublicKey)
					if !writeStatus {
						log.Fatal("Invalid login. Change login credentials")
					}

				} else {
					fmt.Println("=== Invalid session on session.txt file. Issuing new token ===")
					emptyString := ""
					accessToken, _, privateKey, publicKey = authenticate(emptyString, emptyString, emptyString)
				}
			} else {
				fmt.Println("=== Login using session.txt file ===")
				privateKey = sessionData.PrivateKey
				publicKey = sessionData.PublicKey
				accessToken = sessionData.AccessToken
			}
		}
	}

	// Important : Set the accessToken in Environment variable to access in later
	os.Setenv("accessToken", accessToken)

	/* AUTH BLOCK ENDS */

	/*************************************
	// SXT Biscuit
	*************************************/

	// Operation values can include
	/*
		SQL operations:
		DDL: for operating on database objects & schema (e.g. table, view)
		ddl_create : SQL CREATE command
		ddl_alter : SQL ALTER command
		ddl_drop : SQL DROP command

		DML: for performing data manipulation
		dml_insert : SQL INSERT command
		dml_update : SQL UPDATE command
		dml_merge : DQL MERGE command
		dml_delete : SQL DELETE command

		DQL: for performing queries
		dql_select : SQL SELECT command

		Kafka ICM operations:
		kafka_icm_create : ICM create
		kafka_icm_read : ICM read
		kafka_icm_update : ICM update
		kafka_icm_delete : ICM delete

		For wildcards, use *, like
		authorization.SxTBiscuitStruct{Operation: "*", Resource: "eth.TESTTABLE106"}

	*/
	sxtBiscuitCapabilities = append(sxtBiscuitCapabilities, authorization.SxTBiscuitStruct{Operation: "dql_select", Resource: "eth.testtable106"})
	sxtBiscuitCapabilities = append(sxtBiscuitCapabilities, authorization.SxTBiscuitStruct{Operation: "dml_insert", Resource: "eth.testtable106"})
	sxtBiscuitCapabilities = append(sxtBiscuitCapabilities, authorization.SxTBiscuitStruct{Operation: "dml_update", Resource: "eth.testtable106"})
	sxtBiscuitCapabilities = append(sxtBiscuitCapabilities, authorization.SxTBiscuitStruct{Operation: "dml_delete", Resource: "eth.testtable106"})
	sxtBiscuitCapabilities = append(sxtBiscuitCapabilities, authorization.SxTBiscuitStruct{Operation: "ddl_create", Resource: "eth.testtable106"})

	biscuit, _ := authorization.CreateBiscuitToken(sxtBiscuitCapabilities, &privateKey)

	// Set some constant values 
	// Create multiple biscuits
	mb := []string{biscuit}
	originApp := "TEST"


	/*************************************
	// SQL CORE
	*************************************/

	// DDL
	// Create a new schema
	// This wont work as ETH already exists. Write your own
	sqlcore.CreateSchema("CREATE SCHEMA ETH", originApp, mb)

	// DDL
	// Only for create queries
	// For ALTER and DROP, use sqlcore.DDL()
	sqlcore.CreateTable("CREATE TABLE ETH.TESTTABLE106 (id INT PRIMARY KEY, test VARCHAR)", "ALL", originApp, mb, publicKey)



	// DML
	// use the sqlcore.DML to write insert, update, delete, and merge queries
	resources := []string{"ETH.TESTTABLE106"}
	sqlcore.DML("insert into ETH.TESTTABLE106 values(5, 'x5')", originApp, mb, resources)

	// DQL
	// Select operations
	// If rowCount is 0, then fetches all data without limit
	d, e, s := sqlcore.DQL("select * from ETH.TESTTABLE106", originApp, mb, resources, 0)
	if !s {
		log.Println(e)
	} else {
		log.Println(d)
	}

	// DDL
	// Only for ALTER and DROP queries
	// For Create table queries, use sqlcore.CreateTable()
	// For this to work, you will need to provide permissions in biscuit first
	sqlcore.DDL("DROP TABLE ETH.TESTTABLE106", originApp, mb  )

	/*************************************
	// SQL VIEW
	*************************************/

	// Create parameters to pass to view requests
	// Optional param for sqlview.Create(). Not required if viewText param doesnt contain any paramter
	var parametersRequest []sqlview.ParametersRequest
	// parametersRequest = append(parametersRequest, sqlview.ParametersRequest{Name: "count", Type: "integer"})

	// CREATE VIEW
	sqlview.Create("ETH.TESTTABLE106", "testview61", "select count(id) as x from ETH.TESTTABLE106", "Test description", true, parametersRequest)

	// EXECUTE VIEW
	sqlview.Execute("testview61", parametersRequest)

	// UPDATE VIEW
	sqlview.Update("ETH.TESTTABLE106", "testview61", "select count(id) as x from ETH.TESTTABLE106", "Test description", true, parametersRequest)

	// DELETE VIEW
	sqlview.Delete("testview61")

	/*************************************
	// DISCOVERY APIS
	*************************************/

	// List Schemas
	discovery.ListSchemas("ALL", "")

	// List Tables in a given schema
	// Possible scope values -  
	// 1. ALL = all resources, 
	// 2. PUBLIC = non-permissioned tables, 
	// 3. PRIVATE = tables created by the requesting user
	// 4. SUBSCRIPTION =  include all resources the requesting user can see
	// Note: schema and table names are case sensitive. Upper cased
	discovery.ListTables("ETH", "ALL", "")

	// List Columns for a given table in schema
	// Note: schema and table names are case sensitive. Upper cased
	discovery.ListColumns("ETH", "TESTTABLE106")

	// List table index for a given table in schema
	// Note: schema and table names are case sensitive. Upper cased
	discovery.ListTableIndex("ETH", "TESTTABLE106")

	// List table primary key for a given table in schema
	// Note: schema and table names are case sensitive. Upper cased
	discovery.ListTablePrimaryKey("ETH", "TESTTABLE106")

	// List table relations for a schema and scope
	// Note: schema and table names are case sensitive. Upper cased
	discovery.ListTableRelations("ETH", "PRIVATE")

	// List table primary key references for a table and a schema
	// Note: schema and table names are case sensitive. Upper cased
	discovery.ListPrimaryKeyReferences("ETH", "TESTTABLE106", "TEST")

	// List foreign key references for a table, column and a schema
	// Note: schema and table names are case sensitive. Upper cased
	discovery.ListForeignKeyReferences("ETH", "TESTTABLE106", "TEST")

	// List views for a view name for a owner
	// Second parameter can be "", "true" or "false". It accepts a string and not a boolean
	// Note: schema and table names are case sensitive. Upper cased
	discovery.ListViews("SOME_VIEW_NAME", "")

}
