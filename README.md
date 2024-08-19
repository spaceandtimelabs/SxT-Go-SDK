# go-sxt-sdk (v.0.0.8)

Golang SDK for Space and Time Gateway (go version >= 1.18)

## Installation instructions

```sh
go get github.com/spaceandtimelabs/SxT-Go-SDK
```

## Running locally

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
go run main.go -userid=<USERID> -pubkey=<BASE64 STD ENCODED PUBLIC KEY WITHOUT PADDING> -privkey=<BASE64 STD ENCODED PRIVATE KEY WITHOUT PADDING>

// e.g
// pubkey=SiQrMfU+TfRrqqeo/ZDoOwSrHd9zrG1BCU4oDz+4C4Q
// privkey=ys3hQPyfojJOzNymc0eWOKUiogQFGv3G+eeEDUBB8jpKJCsx9T5N9Guqp6j9kOg7BKsd33OsbUEJTigPP7gLhA
```

This will bypass the `userid` and `joincode` mentioned in **.env** file.

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
func Authenticate()(accessToken, refreshToken string, privateKey ed25519.PrivateKey, publicKey ed25519.PublicKey){

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

```go
// Create a new schema
sqlcore.CreateSchema("CREATE SCHEMA ETH")

// Only for create queries
// For ALTER and DROP, use sqlcore.DDL()
sqlcore.CreateTable("ETH.TESTTABLE103", "CREATE TABLE ETH.TESTTABLE103 (id INT PRIMARY KEY, test VARCHAR)", "permissioned", biscuit, publicKey)

// Only for ALTER and DROP queries
// For Create table queries, use sqlcore.CreateTable()
sqlcore.DDL("ETH.TESTTABLE103", "ALTER TABLE ETH.TESTTABLE103 ADD TEST2 VARCHAR", biscuit)

// DML
// use the sqlcore.DML to write insert, update, delete, and merge queries
sqlcore.DML("ETH.TESTTABLE103", "insert into ETH.TESTTABLE103 values(5, 'x5')", biscuit);

// DQL
// Select operations
// If rowCount is 0, then fetches all data without limit
sqlcore.DQL("ETH.TESTTABLE103", "select * from ETH.TESTTABLE103", biscuit, 0);
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

## Configuring a project with SxT SDK

1. Import library

```go
go get github.com/spaceandtimelabs/SxT-Go-SDK
```

2. (Optional) Create a tmp folder in the project root to save user session information (default option). The other options include AWS secrets manager which is not discussed here

```sh
mkdir tmp
```

3. Copy all the `.env.sample` parameters to your environment file

```sh
BASEURL_GENERAL="https://<base_url>/v1" # Space and Time General API Endpoint
BASEURL_DISCOVERY="https://<base_url>/v2"  # Space and Time Discovery API Endpoint
USERID="" # UserID required for authentication and authorization
JOINCODE="" # Space and Time Join Code which can be got from the SxT release team
SCHEME="ed25519"  # The key scheme or algorithm required for key generation
```

4. Integration code

```go
func main() {

	godotenv.Load(".env")

	var biscuits, resources []string

	inputUserID := os.Getenv("USERID")
	pubKey := os.Getenv("PUB_KEY")
	privKey := os.Getenv("PRIV_KEY")

	// Private key to byte array
	pvtKeyBytes, err := base64.StdEncoding.DecodeString(privKey)
	if err != nil {
		log.Println("Private key decoding to []bytes error", err)
	}

	// public key
	// Some languages generate 32-byte private key while some generate 64-byte ones. For such cases, 64-byte pvt key = 32-byte actual private key + 32-byte public key
	pubKeyBytes, err := base64.StdEncoding.DecodeString(pubKey)
	if err != nil {
		log.Println("Public key decoding to []bytes error", err)
	}
	if len(pvtKeyBytes) < 64 {
		pubKeyBytes = append(pubKeyBytes, pubKeyBytes...)
		privKey = base64.StdEncoding.EncodeToString(pvtKeyBytes)
	}

	// Authenticate and get accessToken
	accessToken, _, _, _, err := utils.Authenticate(inputUserID, pubKey, privKey)
	if err != nil {
		log.Println("Authentication error: ", err)
	}

	// Important: Save access token to env
	os.Setenv("accessToken", accessToken)

	// Preparing to call DML
	biscuits = append(biscuits, "actual_biscuit_string")
	resources = append(resources, "schema_name.table_name")

	sqlQuery :=  "INSERT INTO schema_name.tablename VALUES(....)"

	errString, status := sqlcore.DML(sqlQuery, "", biscuits, resources)
	if !status{
		log.Println("Error inserting record to Space and Time: ", errString)
	}

}
```
