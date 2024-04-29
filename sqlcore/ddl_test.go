package sqlcore

import (
	"flag"
	"fmt"
	"path/filepath"
	"testing"

	"github.com/joho/godotenv"
)


func init() {
    if flag.Lookup("test.v") == nil {
        fmt.Println("normal run")
    } else {
        fmt.Println("run under go test")
    }
}
func TestCreateSchema(t *testing.T){

	err := godotenv.Load(filepath.Join("../", ".env.test"))
	if err != nil {
		t.Log("Error: ", err)
	}

	x := []string{}
    errMsg, status := DDL("CREATE SCHEMA ETHEREUM", "TEST", x)
	if !status {
		t.Log(errMsg)
		t.Fail()
	}
    // want := 10

    // if got != want {
    //     t.Errorf("got %q, wanted %q", got, want)
    // }
}