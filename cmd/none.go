package cmd

import (
	"encoding/json"
	"fmt"
	"jwttool/pkg/jwt"
	"os"
)

func NoneAtack(tokenString string, newPayload map[string]interface{}) {
	newToken, err := jwt.NoneAtack(tokenString, newPayload)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error performing none attack: %v\n", err)
		return
	}
	fmt.Println("New token with 'none' algorithm:")
	fmt.Println(newToken)

	parsed, err := jwt.ParseToken(newToken)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing new token payload: %v\n", err)
		return
	}
	jsonbytes, err := json.MarshalIndent(parsed.Payload, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error formatting JSON: %v\n", err)
		return
	}
	fmt.Println("New Payload:")
	fmt.Println(string(jsonbytes))
}
