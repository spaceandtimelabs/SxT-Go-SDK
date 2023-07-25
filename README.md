# go-sxt-sdk (v.0.0.3)

Golang SDK for Space and Time Gateway (go version >= 1.18)

## Installation Instructions

_Note: Before running the code, rename `.env.sample` to `.env` and ensure that your credentials are setup in the `.env` file properly. You will need to obtain a `joinCode` and `endpoint` before you can begin_

```sh
go mod tidy
```

## Features

-   **Sessions**

    The sdk can implement persistent storage in

1. _File based sessions_
2. _AWS Secrets Manager_.

It implements API V2 of Aws SDK (https://github.com/aws/aws-sdk-go-v2). Also access keys, access secrets are retrieved from sharedConfig and sharedCredentials. Read more here https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-files.html

-   **Encryption**

    Support for Ed25519 public key encryption for Biscuit Authorization and Securing data in the platform.

-   **SQL support**

    -   Support for DDL: `creating own schema (namespace), tables`; `altering and deleting tables`
    -   DML: Write all `CRUD` operations
    -   DQL: Any `select` operations
    -   SQL Views

-   **Platform Discovery**
    For fetching meta data from the platform
    -   Namespaces
    -   Tables
    -   Table columns
    -   Table index
    -   Table primary key
    -   Table relations
    -   Primary key references
    -   Foreign key references
    -   List views

## Examples

-   **Running all examples at once**

`main.go` contains a complete running example. To run the file

```go
go run main.go
```

To pass an existing user, use

```go
go run main.go -userid=<USERID> -pubkey=<BASE64 STD ENCODED PUBLIC KEY> -privkey=<BASE64 STD ENCODED PRIVATE KEY>

// e.g
// pubkey=SiQrMfU+TfRrqqeo/ZDoOwSrHd9zrG1BCU4oDz+4C4Q=
// privkey=ys3hQPyfojJOzNymc0eWOKUiogQFGv3G+eeEDUBB8jpKJCsx9T5N9Guqp6j9kOg7BKsd33OsbUEJTigPP7gLhA=

// The keys are provided for example and will not work
```

This will bypass the `userid` and `joincode` mentioned in **.env** file.

**Note**
_SxT libraries may generate a 32-byte or a 64-byte base64 encoded private key. A 64-byte is a combination of 32-byte secret and 32-byte public key_

For details on running the `main.go` file, use

```go
go run main.go -h
```

-   **Authentication**

It is very important to save your **private key** used in authentication and biscuit generation. Else you will not have access to the user and tables created using the key.

The generated `AccessToken` is valid for 25 minutes and the `refreshToken` for 30 minutes.

```go
// New Authentication.
// Generates new accessToken, refreshToken, privateKey, and publicKey
func authenticate()(accessToken, refreshToken string, privateKey ed25519.PrivateKey, publicKey ed25519.PublicKey){

    // Read userId, joinCode from .env file
	userId, _ := helpers.ReadUserId()
	joinCode, _ := helpers.ReadJoinCode()

	var pubkey ed25519.PublicKey
	var privkey ed25519.PrivateKey

	var authCodeStruct authentication.AuthCodeStruct
	var tokenStruct authentication.TokenStruct

	sessionStruct, sessionStatus := storage.FileReadSession(userId)

	if !sessionStatus {
		// Generate Private and Public keys
		pubkey, privkey = helpers.CreateKey()
	} else {
		pubkey = sessionStruct.PublicKey
		privkey = sessionStruct.PrivateKey
	}


	// Get auth code
	authCode := authentication.GenerateAuthCode(userId, joinCode)
	json.Unmarshal([]byte(authCode), &authCodeStruct)

	// Get Keys
	encodedSignature,  base64PublicKey := authentication.GenerateKeys(authCodeStruct.AuthCode, pubkey, privkey)

	// Get Token
	tokenJson := authentication.GenerateToken(userId, authCodeStruct.AuthCode, encodedSignature, base64PublicKey)
	json.Unmarshal([]byte(tokenJson), &tokenStruct)

	// Store session data to a local file

	writeStatus := storage.FileWriteSession(userId, tokenStruct.AccessToken, tokenStruct.RefreshToken, privkey, pubkey)
	if !writeStatus{
		log.Fatal("Invalid login. Change login credentials")
	}

	return tokenStruct.AccessToken, tokenStruct.RefreshToken, privkey, publicKey
}

```

-   **Generating Biscuits**

You can create multiple biscuit tokens for a table allowing you provide different access levels for users. For the list of all the capabilities, [visit https://docs.spaceandtime.io/docs/biscuit-authorization](https://docs.spaceandtime.io/docs/biscuit-authorization)

Sample biscuit generation with permissions for `select query`, `insert query`, `update query`, `delete query`, `create table`

```go
var sxtBiscuitCapabilities []authorization.SxTBiscuitStruct

// Add biscuit capabilities
sxtBiscuitCapabilities = append(sxtBiscuitCapabilities, authorization.SxTBiscuitStruct{Operation: "dql_select", Resource: "eth.testtable103"})
sxtBiscuitCapabilities = append(sxtBiscuitCapabilities, authorization.SxTBiscuitStruct{Operation: "dml_insert", Resource: "eth.testtable103"})
sxtBiscuitCapabilities = append(sxtBiscuitCapabilities, authorization.SxTBiscuitStruct{Operation: "dml_update", Resource: "eth.testtable103"})
sxtBiscuitCapabilities = append(sxtBiscuitCapabilities, authorization.SxTBiscuitStruct{Operation: "dml_delete", Resource: "eth.testtable103"})
sxtBiscuitCapabilities = append(sxtBiscuitCapabilities, authorization.SxTBiscuitStruct{Operation: "ddl_create", Resource: "eth.testtable103"})

// Generate the biscuit token
biscuit, _ := authorization.CreateBiscuitToken(sxtBiscuitCapabilities, &privateKey)

```

-   **DDL, DML & DQL**

    **Note**:

To generate a new **schema**, `ddl_create` permission is needed

For multi biscuit support and setting `originApp` for identifying logs

```go
// Create multiple biscuits
mb := []string{biscuit}
originApp := "TEST"
```

#### Queries

```go
// Create a new schema
sqlcore.CreateSchema("CREATE SCHEMA ETH")

// Only for create queries
// For DROP use sqlcore.DDL()
sqlcore.CreateTable("CREATE TABLE ETH.TESTTABLE103 (id INT PRIMARY KEY, test VARCHAR)", "permissioned", biscuit, originApp, mb, publicKey)

// Only for DROP queries
// For Create table queries, use sqlcore.CreateTable()
sqlcore.DDL("DROP TABLE ETH.TESTTABLE103", biscuit, originApp, mb )

// DML
// use the sqlcore.DML to write insert, update, delete, and merge queries
sqlcore.DML("ETH.TESTTABLE103", "insert into ETH.TESTTABLE103 values(5, 'x5')", biscuit, originApp, mb);

// DQL
// Select operations
// If rowCount is 0, then fetches all data without limit
sqlcore.DQL("ETH.TESTTABLE103", "select * from ETH.TESTTABLE103", biscuit, originApp, mb, 0);
```

-   **DISCOVERY**

Discovery calls need a user to be logged in

```go

// List Namespaces
discovery.ListNamespaces()

// List Tables in a given namespace
// Possible scope values -  ALL = all resources, PUBLIC = non-permissioned tables, PRIVATE = tables created by the requesting user
discovery.ListTables("ETH", "ALL")

// List Columns for a given table in namespace
discovery.ListColumns("ETH", "TESTTABLE103")

// List table index for a given table in namespace
discovery.ListTableIndex("ETH", "TESTTABLE103")


// List table primary key for a given table in namespace
discovery.ListTablePrimaryKey("ETH", "TESTTABLE103")

// List table relations for a namespace and scope
// Possible scope values -  ALL = all resources, PUBLIC = non-permissioned tables, PRIVATE = tables created by the requesting user
discovery.ListTableRelations("ETH", "PRIVATE")

// List table primary key references for a table and a namespace
discovery.ListPrimaryKeyReferences("ETH", "TESTTABLE103", "TEST")

// List foreign key references for a table, column and a namespace
discovery.ListForeignKeyReferences("ETH", "TESTTABLE103", "TEST")
```

-   **Storage**

For AWS and File storage, the following methods are available

```go

// File
storage.FileWriteSession(userId, tokenStruct.AccessToken, tokenStruct.RefreshToken, privkey, pubkey)

storage.FileReadSession(userId)

storage.FileUpdateSession(userId, accessToken, refreshToken, privateKey, publicKey)

// AWS
storage.AwsWriteSession(userId, tokenStruct.AccessToken, tokenStruct.RefreshToken, privkey, pubkey)

storage.AwsReadSession(userId)

storage.AwsUpdateSession(userId, accessToken, refreshToken, privateKey, publicKey)
```
