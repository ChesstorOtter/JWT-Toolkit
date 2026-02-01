package cmd

import (
	"fmt"
	"jwttool/pkg/jwt"
	"os"
)

func VerifyToken(tokenString string, secret string) {
	valid, err := jwt.VerifyHS256(tokenString, []byte(secret))
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error verifying token:", err)
		return
	}
	if valid {
		fmt.Println("Token is valid with secret.", secret)
	} else {
		fmt.Println("Token is invalid with secret.", secret)
	}
}
