package jwt

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
)

func VerifyHS256(tokenString string, secret []byte) (bool, error) {
	parts := strings.Split(tokenString, ".")
	if len(parts) != 3 {
		return false, fmt.Errorf("Invalid token format")
	}

	signature, err := base64.RawURLEncoding.DecodeString(parts[2])
	if err != nil {
		return false, fmt.Errorf("Invalid signature")
	}

	mac := hmac.New(sha256.New, secret)
	mac.Write([]byte(parts[0] + "." + parts[1]))
	expectedSignature := mac.Sum(nil)

	return hmac.Equal(signature, expectedSignature), nil
}

func SignHS256(header, payload map[string]interface{}, secret []byte) (string, error) {
	headerBytes, err := json.Marshal(header)
	if err != nil {
		return "", fmt.Errorf("Failed to marshal header: %v", err)
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("Failed to marshal payload: %v", err)
	}
	headerEncoded := base64.RawURLEncoding.EncodeToString(headerBytes)
	payloadEncoded := base64.RawURLEncoding.EncodeToString(payloadBytes)
	signingInput := fmt.Sprintf("%s.%s", headerEncoded, payloadEncoded)
	mac := hmac.New(sha256.New, secret)
	mac.Write([]byte(signingInput))
	signature := mac.Sum(nil)
	signatureEncoded := base64.RawURLEncoding.EncodeToString(signature)
	tokenString := fmt.Sprintf("%s.%s.%s", headerEncoded, payloadEncoded, signatureEncoded)
	return tokenString, nil
}
