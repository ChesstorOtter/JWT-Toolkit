package jwt

import (
	"encoding/base64"
	"encoding/json"
	//"errors"
	"fmt"
	"strings"
)

type Token struct {
	Header    map[string]interface{}
	Payload   map[string]interface{}
	Signature string
	Raw       string
}

func ParseToken(tokenString string) (*Token, error) {
	parts := strings.Split(tokenString, ".")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid token format")
	}

	headerBytes, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return nil, fmt.Errorf("failed to decode header: %v", err)
	}

	var header map[string]interface{}
	if err := json.Unmarshal(headerBytes, &header); err != nil {
		return nil, fmt.Errorf("failed to unmarshal header: %v", err)
	}

	payloadBytes, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, fmt.Errorf("failed to decode payload: %v", err)
	}

	var payload map[string]interface{}
	if err := json.Unmarshal(payloadBytes, &payload); err != nil {
		return nil, fmt.Errorf("failed to unmarshal payload: %v", err)
	}

	return &Token{
		Header:    header,
		Payload:   payload,
		Signature: parts[2],
		Raw:       tokenString,
	}, nil

}
