package helpers

import (
	"regexp"
	"strings"
)

// Performs a regex check to validate if the input is upper case only with numbers
func CheckUpperCase(input string) (msg string, status bool){

	parts := strings.Split(input, ".")

	if len(parts) > 1 {
		return "invalid resource name", false
	}

	m, _ := regexp.Match("^[A-Z_][A-Z0-9_]*$", []byte(input))
	if !m {
		return "input should be upper case", m
	}

	return "ok", true
}


// Performs a regex check to validate if the resourceID is upper case only with numbers
func CheckUpperCaseResource(input string) (msg string, status bool){
	parts := strings.Split(input, ".")

	if len(parts) == 1 {
		return "please prefix 'PUBLIC.' to resourceID ", false
	} else {
		for _, part := range parts{
			m, _ := regexp.Match("^[A-Z_][A-Z0-9_]*$", []byte(part))
			if !m {
				return "input should be upper case", m
			}
		}
	}

	return "ok", true
}