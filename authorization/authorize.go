package authorization

import (
	"crypto/ed25519"
	"encoding/base64"

	"github.com/biscuit-auth/biscuit-go/v2"
	"github.com/biscuit-auth/biscuit-go/v2/parser"
)

type SxTBiscuitStruct struct{

	Operation string
	Resource string
}


// Create Biscuit Token
func CreateBiscuitToken(capabilities []SxTBiscuitStruct, root *ed25519.PrivateKey) (biscuitToken string, status bool) {
	
	var capabilityString string

	builder := biscuit.NewBuilder(*root)
	
	for _, capability := range capabilities {
		capabilityString = `sxt:capability("` + capability.Operation + `","` + capability.Resource + `")`

		fact, err := parser.FromStringFact(capabilityString)
		if err != nil {
			return "", false
		}
		err = builder.AddAuthorityFact(fact)
		if err != nil {
			return "", false
		}
    }

	token, err := builder.Build()
	if err != nil {
		return "", false
	}

	tokenSerialized, err := token.Serialize()
	if err != nil {
		return "", false
	}

	return base64.URLEncoding.EncodeToString(tokenSerialized), true
}