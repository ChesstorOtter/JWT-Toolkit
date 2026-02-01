package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
)

func main() {
	// Генерируем RSA ключи
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error generating private key: %v\n", err)
		return
	}

	// Сохраняем приватный ключ
	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	privateKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyBytes,
	})
	err = os.WriteFile("private_key.pem", privateKeyPEM, 0600)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing private key: %v\n", err)
		return
	}

	// Сохраняем публичный ключ
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error marshaling public key: %v\n", err)
		return
	}
	publicKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKeyBytes,
	})
	err = os.WriteFile("public_key.pem", publicKeyPEM, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing public key: %v\n", err)
		return
	}

	fmt.Println("Keys generated successfully!")
	fmt.Println("  private_key.pem - saved")
	fmt.Println("  public_key.pem - saved")
}
