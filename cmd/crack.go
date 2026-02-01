package cmd

import (
	"fmt"
	"jwttool/pkg/jwt"
)

func CrackToken(tokenString string, wordlistPath string) {
	fmt.Println("Starting token cracking...")
	jwt.CrackHS256(tokenString, wordlistPath)
}
