package helpers

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Read User Id from Environment
func ReadUserId() (value string, ok bool){
	
	value, ok = os.LookupEnv("USERID")
	if !ok {
		envFile, _ := godotenv.Read(".env")
		value = envFile["USERID"]

		if value == ""{
			log.Fatal("USERID not set in environment")
		}
	}

	return value, true
}

// Read Join Code from Environment 
func ReadJoinCode() (value string, ok bool){
	value, ok = os.LookupEnv("JOINCODE")
	if !ok {
		envFile, _ := godotenv.Read(".env")
		value = envFile["JOINCODE"]

		if value == ""{
			log.Fatal("JOINCODE not set in environment")
		}
	}
	
	return value, true
}

// Read API End Point Discovery from Environment 
func ReadEndPointDiscovery() (value string, ok bool){
	value, ok = os.LookupEnv("BASEURL_DISCOVERY")
	if !ok {
		envFile, _ := godotenv.Read(".env")
		value = envFile["BASEURL_DISCOVERY"]

		if value == ""{
			log.Fatal("BASEURL_DISCOVERY not set in environment")
		}
	}
	return value, true
}

// Read API End Point Others in General from Environment 
func ReadEndPointGeneral() (value string, ok bool){

	value, ok = os.LookupEnv("BASEURL_GENERAL")
	if !ok {
		envFile, _ := godotenv.Read(".env")
		value = envFile["BASEURL_GENERAL"]

		if value == ""{
			log.Fatal("BASEURL_GENERAL not set in environment")
		}
	}

	return value, true
}

// Read Scheme from Environment 
func ReadScheme() (value string, ok bool){
	value, ok = os.LookupEnv("SCHEME")
	if !ok {
		envFile, _ := godotenv.Read(".env")
		value = envFile["SCHEME"]

		if value == ""{
			log.Fatal("SCHEME not set in environment")
		}
	}
	return value, true
}