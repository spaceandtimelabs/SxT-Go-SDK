package utils

import (
	"crypto/ed25519"
	"encoding/base64"
	"encoding/json"
	"errors"
	"log"

	"github.com/spaceandtimelabs/SxT-Go-SDK/authentication"
	"github.com/spaceandtimelabs/SxT-Go-SDK/authorization"
	"github.com/spaceandtimelabs/SxT-Go-SDK/discovery"
	"github.com/spaceandtimelabs/SxT-Go-SDK/helpers"
	"github.com/spaceandtimelabs/SxT-Go-SDK/sqlcore"
	"github.com/spaceandtimelabs/SxT-Go-SDK/storage"
)

// New Authentication.
// This method generates new accessToken, refreshToken, privateKey, and publicKey
func Authenticate(inputUserId, inputPublicKey, inputPrivateKey string) ( accessToken, refreshToken string, privateKey ed25519.PrivateKey, publicKey ed25519.PublicKey, err error) {

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
			return "", "", nil, nil, errors.New("base64 std encoded public key expected")
		}

		privkey, e = base64.StdEncoding.DecodeString(inputPrivateKey)
		if e != nil {
			log.Fatal("base64 std encoded private key expected")
			return  "", "", nil, nil, errors.New("base64 std encoded private key expected")
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
		// log.Fatal("Invalid login. Change login credentials")
		return "", "", nil, nil, errors.New("invalid login. Change login credentials")
	}

	// Logout
	// authentication.Logout()
	return  tokenStruct.AccessToken, tokenStruct.RefreshToken, privkey, publicKey, nil
}

// SQL APIs
func SQLAPIs(privateKey ed25519.PrivateKey, publicKey ed25519.PublicKey) (err error){

	var sxtBiscuitCapabilities []authorization.SxTBiscuitStruct

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
	errMsg, status := sqlcore.CreateSchema("CREATE SCHEMA ETH", originApp, mb)
	if !status {
		return errors.New(errMsg)
	}

	// DDL
	// Only for create queries
	// For ALTER and DROP, use sqlcore.DDL()
	sqlcore.CreateTable("CREATE TABLE ETH.TESTTABLE106 (id INT PRIMARY KEY, test VARCHAR)", "ALL", originApp, mb, publicKey)
	if !status {
		return errors.New(errMsg)
	}

	// DML
	// use the sqlcore.DML to write insert, update, delete, and merge queries
	resources := []string{"ETH.TESTTABLE106"}
	sqlcore.DML("insert into ETH.TESTTABLE106 values(5, 'x5')", originApp, mb, resources)
	if !status {
		return errors.New(errMsg)
	}

	// DQL
	// Select operations
	// If rowCount is 0, then fetches all data without limit
	_, errMsg, status = sqlcore.DQL("select * from ETH.TESTTABLE106", originApp, mb, resources, 0)
	if !status {
		return errors.New(errMsg)
	}

	// DDL
	// Only for ALTER and DROP queries
	// For Create table queries, use sqlcore.CreateTable()
	// For this to work, you will need to provide permissions in biscuit first
	errMsg, status = sqlcore.DDL("DROP TABLE ETH.TESTTABLE106", originApp, mb  )
	if !status {
		return errors.New(errMsg)
	}

	return nil
}


// Discovery APIs
func DiscoveryAPIs()(err error){

	/*************************************
	// DISCOVERY APIS
	*************************************/

	// List Schemas
	_, errMsg, status := discovery.ListSchemas("ALL", "")
	if !status {
		return errors.New(errMsg)
	}

	// List Tables in a given schema
	// Possible scope values -  
	// 1. ALL = all resources, 
	// 2. PUBLIC = non-permissioned tables, 
	// 3. PRIVATE = tables created by the requesting user
	// 4. SUBSCRIPTION =  include all resources the requesting user can see
	// Note: schema and table names are case sensitive. Upper cased
	_, errMsg, status = discovery.ListTables("ETH", "ALL", "")
	if !status {
		return errors.New(errMsg)
	}

	// List Columns for a given table in schema
	// Note: schema and table names are case sensitive. Upper cased
	_, errMsg, status = discovery.ListColumns("ETH", "TESTTABLE106")
	if !status {
		return errors.New(errMsg)
	}

	// List table index for a given table in schema
	// Note: schema and table names are case sensitive. Upper cased
	_, errMsg, status = discovery.ListTableIndex("ETH", "TESTTABLE106")
	if !status {
		return errors.New(errMsg)
	}

	// List table primary key for a given table in schema
	// Note: schema and table names are case sensitive. Upper cased
	_, errMsg, status = discovery.ListTablePrimaryKey("ETH", "TESTTABLE106")
	if !status {
		return errors.New(errMsg)
	}

	// List table relations for a schema and scope
	// Note: schema and table names are case sensitive. Upper cased
	_, errMsg, status = discovery.ListTableRelations("ETH", "PRIVATE")
	if !status {
		return errors.New(errMsg)
	}

	// List table primary key references for a table and a schema
	// Note: schema and table names are case sensitive. Upper cased
	_, errMsg, status = discovery.ListPrimaryKeyReferences("ETH", "TESTTABLE106", "TEST")
	if !status {
		return errors.New(errMsg)
	}

	// List foreign key references for a table, column and a schema
	// Note: schema and table names are case sensitive. Upper cased
	_, errMsg, status = discovery.ListForeignKeyReferences("ETH", "TESTTABLE106", "TEST")
	if !status {
		return errors.New(errMsg)
	}

	// List views for a view name for a owner
	// Second parameter can be "", "true" or "false". It accepts a string and not a boolean
	// Note: schema and table names are case sensitive. Upper cased
	_, errMsg, status = discovery.ListViews("SOME_VIEW_NAME", "")
	if !status {
		return errors.New(errMsg)
	}
	
	return nil
}