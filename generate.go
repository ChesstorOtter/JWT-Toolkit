package main

import (
	"fmt"
	"jwttool/pkg/jwt"
	"os"
)

func main() {
	secret := []byte("secret")

	header := map[string]interface{}{
		"alg": "HS256",
		"typ": "JWT",
	}

	payload := map[string]interface{}{
		"user": "bob",
		"role": "admin",
		"iat":  1700000000,
	}

	token, err := jwt.SignHS256(header, payload, secret)
	if err != nil {
		fmt.Println("Error generating token:", err)
		return
	}

	wordlist := []string{
		"password", "123456", "admin", "test", "qwerty",
		"secret", "letmein", "welcome", "login", "pass",
		"secret123", "password123", "admin123", "12345678", "1234567",
		"123123", "abc123", "monkey", "1234", "dragon",
		"master", "sunshine", "princess", "qwertyuiop", "solo",
		"passw0rd", "starwars", "shadow", "michael", "football",
		"batman", "superman", "ironman", "spiderman", "thanos",
		"admin@123", "root", "toor", "password@123", "test@123",
		"user", "guest", "demo", "test123", "admin@",
		"pwd", "pass", "secret@123", "keyboard", "trustno1",
		"twitter", "facebook", "instagram", "github", "gitlab",
		"mongodb", "postgres", "mysql", "redis", "docker",
		"kubernetes", "jenkins", "gitlab", "github", "bitbucket",
		"spring", "hibernate", "hibernate", "servlet", "jsp",
		"node", "express", "react", "angular", "vue",
		"golang", "rust", "python", "java", "csharp",
	}

	file, err := os.Create("wordlist.txt")
	if err != nil {
		fmt.Println("Error creating wordlist.txt:", err)
		return
	}
	defer file.Close()

	for _, word := range wordlist {
		file.WriteString(word + "\n")
	}

	fmt.Println("=== JWT Token Generation ===")
	fmt.Println("\nGenerated token:")
	fmt.Println(token)
	fmt.Println("\nWordlist saved to: wordlist.txt")
	fmt.Println("\n=== Testing Commands ===")
	fmt.Printf("Crack token:\n  go run main.go crack \"%s\" \"wordlist.txt\"\n", token)
}
