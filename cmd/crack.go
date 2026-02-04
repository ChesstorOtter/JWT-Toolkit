package cmd

import (
	"fmt"
	"jwttool/pkg/jwt"
)

func CrackToken(tokenString string, wordlistPath string) {
	fmt.Println("Starting token cracking...")
	secret, err := jwt.CrackHS256(tokenString, wordlistPath)
	if err != nil {
		fmt.Printf("Crack finished: %v\n", err)
		return
	}
	fmt.Printf("Crack finished: secret=%q\n", secret)
}
