package cmd

import (
	"encoding/json"
	"fmt"
	"jwttool/pkg/jwt"
	"os"
)

func DecodeToken(tokenString string) {
	token, err := jwt.ParseToken(tokenString)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error decoding token: %v\n", err)
		return
	}

	fmt.Println("Header:")
	fmt.Printf("%v\n", token.Header)
	fmt.Println("Payload:")
	fmt.Printf("%v\n", token.Payload)
	fmt.Println("Signature:")
	fmt.Printf("%s\n", token.Signature)
}

func PrintDecodedToken(data map[string]interface{}) {
	jsonbytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error formatting JSON: %v\n", err)
		return
	}
	fmt.Println(string(jsonbytes))
}
