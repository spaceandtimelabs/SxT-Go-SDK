package helpers

import (
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

// Read User Id from Environment
func ReadUserId() (value string, ok bool){
	err := godotenv.Load(filepath.Join("./", ".env"))
	if err != nil {
		value = os.Getenv("USERID")
		if value == "" {
			log.Fatal("USERID not set in environment")
		}
	}  else {
		value, ok = os.LookupEnv("USERID")

		if !ok {
			log.Fatal("USERID not set in environment")
		}
	}
	
	return value, true
}

// Read Join Code from Environment 
func ReadJoinCode() (value string, ok bool){

	err := godotenv.Load(filepath.Join("./", ".env"))
	if err != nil {
		value = os.Getenv("JOINCODE")
		if value == "" {
			log.Println("JOINCODE not set in environment")
		}
	}  else {
		value, ok = os.LookupEnv("JOINCODE")

		if !ok {
			log.Println("JOINCODE not set in environment")
		}
	}
	
	return value, true
}

// Read API End Point Discovery from Environment 
func ReadEndPointDiscovery() (value string, ok bool){

	err := godotenv.Load(filepath.Join("./", ".env"))
	if err != nil {
		value = os.Getenv("BASEURL_DISCOVERY")
		if value == "" {
			log.Fatal("Discovery BASEURL not set in environment")
		}
	} else {
		value, ok = os.LookupEnv("BASEURL_DISCOVERY")

		if !ok {
			log.Fatal("Discovery BASEURL not set in environment")
		}
	}

	return value, true
}

// Read API End Point Others in General from Environment 
func ReadEndPointGeneral() (value string, ok bool){

	err := godotenv.Load(filepath.Join("./", ".env"))
	if err != nil {
		value = os.Getenv("BASEURL_GENERAL")
		if value == "" {
			log.Fatal("General BASEURL not set in environment")
		}
	} else {
		value, ok = os.LookupEnv("BASEURL_GENERAL")

		if !ok {
			log.Fatal("General BASEURL not set in environment")
		}
	}

	return value, true
}

// Read Scheme from Environment 
func ReadScheme() (value string, ok bool){
	value = "ed25519"
	return value, true
}