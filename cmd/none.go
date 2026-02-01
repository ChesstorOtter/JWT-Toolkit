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

	if payload, ok := newPayload["payload"]; ok {
		fmt.Println("New Payload:")
		jsonbytes, err := json.MarshalIndent(payload, "", "  ")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error formatting JSON: %v\n", err)
			return
		}
		fmt.Println(string(jsonbytes))
	}
}
