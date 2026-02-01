package cmd

import (
	"fmt"
	"jwttool/pkg/jwt"
	"os"
)

func ConfusionAtack(tokenString string, newPayload map[string]interface{}, privateKeyPath string) {
	privateKeyPEM, err := os.ReadFile(privateKeyPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading private key file: %v\n", err)
		return
	}

	newToken, err := jwt.ConfusionAtack(tokenString, newPayload, string(privateKeyPEM))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error performing confusion attack: %v\n", err)
		return
	}
	fmt.Println("New token with 'RS256' algorithm:")
	fmt.Println(newToken)
}
