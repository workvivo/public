package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"log"
	"math/big"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func main() {
	privateKey, publicKey := loadOrGenerateKeyPair()

	// === USER NEEDS TO CONFIGURE THESE SETTINGS HERE ===

	orgId := "165"        // This is the organisation ID
	appWorkvivo := "3981" // APP ID from your Workvivo app with unwiredotp.*.write permission
	appAud := "unwired.workvivo.red"
	apiURL := "https://api-gateway.workvivo.red/v1/unwired/users/otp" //In case it is EU Production HOST should be api.workvivo.com
	postData := `{"email": "test@nomail"}`

	// === KEYPAIR GENERATION OR LOADING ===

	/*
	* You can generate your own RSA keypair manually using OpenSSL:
	*
	* # Generate a 4096-bit private key
	* openssl genpkey -algorithm RSA -out private.pem -pkeyopt rsa_keygen_bits:4096
	*
	* # Extract the public key from the private key
	* openssl rsa -in private.pem -pubout -out public.pem
	*
	* Place both files (private.pem and public.pem) in the same directory as this script.
	* The script will use them automatically instead of generating new keys.
	 */

	// Create JWT
	payload := jwt.MapClaims{
		"jti":         randomHex(32),
		"iss":         "org1",
		"sub":         "app",
		"workvivo_id": appWorkvivo,
		"aud":         appAud,
		"iat":         time.Now().Unix(),
		"nbf":         time.Now().Unix(),
		"exp":         time.Now().Add(10 * time.Minute).Unix(),
		"state":       randomHex(32),
	}

	kid := computeKID(publicKey)
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, payload)
	token.Header["kid"] = kid
	jwtString, err := token.SignedString(privateKey)
	check(err)

	// JWKS
	jwks := createJWKS(publicKey, kid)

	fmt.Println("Token (JWT):\n" + jwtString + "\n")
	fmt.Println("KeyID:\n" + kid + "\n")
	fmt.Println("Public Key (PEM):")
	printPublicKeyPEM(publicKey)
	fmt.Println("\nJWKS:")
	jwksJSON, _ := json.MarshalIndent(jwks, "", "  ")
	fmt.Println(string(jwksJSON))

	fmt.Println("\nPress Enter to confirm that JWKS are updated in your organisation here: 'https://HOST/admin/developers/apps/manage':")
	fmt.Scanln()

	// Perform HTTP request

	req, err := http.NewRequest("POST", apiURL, strings.NewReader(postData))
	check(err)
	req.Header.Set("Workvivo-Id", orgId)
	req.Header.Set("x-workvivo-jwt", jwtString)
	req.Header.Set("x-workvivo-jwt-keyid", kid)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	check(err)
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	check(err)
	fmt.Println("Response:\n" + string(body))
	fmt.Printf("HTTP Code:\n%d\n", resp.StatusCode)
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func randomHex(n int) string {
	rb := make([]byte, n)
	_, err := rand.Read(rb)
	check(err)
	return fmt.Sprintf("%x", rb)
}

func computeKID(pub *rsa.PublicKey) string {
	nBytes := pub.N.Bytes()
	eBytes := big.NewInt(int64(pub.E)).Bytes()
	jwk := map[string]string{
		"e":   base64url(eBytes),
		"kty": "RSA",
		"n":   base64url(nBytes),
	}
	jwkJSON, _ := json.Marshal(jwk)
	hash := sha256.Sum256(jwkJSON)
	return base64url(hash[:])
}

func base64url(b []byte) string {
	return strings.TrimRight(base64.URLEncoding.EncodeToString(b), "=")
}

func createJWKS(pub *rsa.PublicKey, kid string) map[string][]map[string]string {
	return map[string][]map[string]string{
		"keys": {
			{
				"kty": "RSA",
				"alg": "RS256",
				"use": "sig",
				"kid": kid,
				"n":   base64url(pub.N.Bytes()),
				"e":   base64url(big.NewInt(int64(pub.E)).Bytes()),
			},
		},
	}
}

func printPublicKeyPEM(pub *rsa.PublicKey) {
	pubBytes, err := x509.MarshalPKIXPublicKey(pub)
	check(err)
	pem.Encode(os.Stdout, &pem.Block{Type: "PUBLIC KEY", Bytes: pubBytes})
}

func loadOrGenerateKeyPair() (*rsa.PrivateKey, *rsa.PublicKey) {
	privKeyPath := "private.pem"
	pubKeyPath := "public.pem"

	if _, err := os.Stat(privKeyPath); os.IsNotExist(err) {
		fmt.Println("Private key not found. Generating new key pair...")

		privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
		check(err)

		privBytes, err := x509.MarshalPKCS8PrivateKey(privateKey)
		check(err)

		privPem := pem.EncodeToMemory(&pem.Block{
			Type:  "PRIVATE KEY",
			Bytes: privBytes,
		})
		err = os.WriteFile(privKeyPath, privPem, 0600)
		check(err)

		pubBytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
		check(err)
		pubPem := pem.EncodeToMemory(&pem.Block{
			Type:  "PUBLIC KEY",
			Bytes: pubBytes,
		})
		err = os.WriteFile(pubKeyPath, pubPem, 0644)
		check(err)
		return privateKey, &privateKey.PublicKey
	}

	privPem, err := os.ReadFile(privKeyPath)
	check(err)
	block, _ := pem.Decode(privPem)
	if block == nil || block.Type != "PRIVATE KEY" {
		log.Fatal("Invalid private key format")
	}
	parsedKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	check(err)
	privateKey, ok := parsedKey.(*rsa.PrivateKey)
	if !ok {
		log.Fatal("Not an RSA private key")
	}
	return privateKey, &privateKey.PublicKey
}
