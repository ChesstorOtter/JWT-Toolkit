package main

import (
	"encoding/json"
	"fmt"
	"jwttool/cmd"
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: jwttool <command> [args...]")
		fmt.Println("Commands:")
		fmt.Println("  decode <token>")
		fmt.Println("  verify <token> <secret>")
		fmt.Println("  crack <token> <wordlist>")
		fmt.Println("  none <token> <newPayloadJSON>")
		return
	}

	command := os.Args[1]

	switch command {
	case "verify":
		if len(os.Args) < 4 {
			fmt.Println("Usage: jwttool verify <token> <secret>")
			return
		}
		tokenString := os.Args[2]
		secret := os.Args[3]
		cmd.VerifyToken(tokenString, secret)

	case "decode":
		if len(os.Args) < 3 {
			fmt.Println("Usage: jwttool decode <token>")
			return
		}
		tokenString := os.Args[2]
		cmd.DecodeToken(tokenString)

	case "crack":
		if len(os.Args) < 4 {
			fmt.Println("Usage: jwttool crack <token> <wordlist>")
			return
		}
		tokenString := os.Args[2]
		wordlistPath := os.Args[3]
		cmd.CrackToken(tokenString, wordlistPath)

	case "none":
		if len(os.Args) < 4 {
			fmt.Println("Usage: jwttool none <token> <newPayloadJSON>")
			return
		}
		tokenString := os.Args[2]
		payloadJSON := strings.Join(os.Args[3:], " ")

		payloadJSON = strings.TrimPrefix(payloadJSON, "\ufeff")
		payloadJSON = strings.TrimSpace(payloadJSON)

		var newPayload map[string]interface{}
		err := json.Unmarshal([]byte(payloadJSON), &newPayload)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing new payload JSON: %v\n", err)
			fmt.Fprintf(os.Stderr, "Received: %q\n", payloadJSON)
			return
		}
		cmd.NoneAtack(tokenString, newPayload)

	case "confusion":
		if len(os.Args) < 5 {
			fmt.Println("Usage: jwttool confusion <token> <newPayloadJSON> <publicKeyPath>")
			return
		}
		tokenString := os.Args[2]
		payloadJSON := strings.Join(os.Args[3:len(os.Args)-1], " ")
		publicKeyPath := os.Args[len(os.Args)-1]

		payloadJSON = strings.TrimPrefix(payloadJSON, "\ufeff")
		payloadJSON = strings.TrimSpace(payloadJSON)

		var newPayload map[string]interface{}
		err := json.Unmarshal([]byte(payloadJSON), &newPayload)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing new payload JSON: %v\n", err)
			return
		}
		cmd.ConfusionAtack(tokenString, newPayload, publicKeyPath)

	case "attack-none":
		if len(os.Args) < 5 {
			fmt.Println("Usage: jwttool attack-none <token> <target> <endpoint> [proxy] [insecure:true|false]")
			return
		}
		tokenString := os.Args[2]
		target := os.Args[3]
		endpoint := os.Args[4]
		proxy := ""
		insecure := false
		if len(os.Args) >= 6 {
			proxy = os.Args[5]
		}
		if len(os.Args) >= 7 {
			insecure = os.Args[6] == "true"
		}
		cmd.AttackNone(tokenString, target, endpoint, proxy, insecure)

	case "attack-crack":
		if len(os.Args) < 6 {
			fmt.Println("Usage: jwttool attack-crack <token> <wordlist> <target> <endpoint> [proxy] [insecure:true|false]")
			return
		}
		tokenString := os.Args[2]
		wordlistPath := os.Args[3]
		target := os.Args[4]
		endpoint := os.Args[5]
		proxy := ""
		insecure := false
		if len(os.Args) >= 7 {
			proxy = os.Args[6]
		}
		if len(os.Args) >= 8 {
			insecure = os.Args[7] == "true"
		}
		cmd.AttackCrack(tokenString, wordlistPath, target, endpoint, proxy, insecure)

	default:
		fmt.Println("Unknown command:", command)
		fmt.Println("Available commands: decode, verify, crack, none, confusion, attack-none, attack-crack")
	}
}
