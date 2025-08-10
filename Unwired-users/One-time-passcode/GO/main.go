package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTPayload struct {
	Aud   string `json:"aud"`
	AppID string `json:"app_id"`
	Sub   string `json:"sub"`
	Iss   string `json:"iss"`
	Nbf   int64  `json:"nbf"`
	Iat   int64  `json:"iat"`
	Exp   int64  `json:"exp"`
	Jti   string `json:"jti"`
	State string `json:"state"`
}

func main() {
	privateKey, publicKey := loadKeyPair()

	/*
		=== USER NEEDS TO CONFIGURE THESE SETTINGS HERE ===
	*/
	/*
	   The Organisation ID is the ID of the Workvivo organisation.
	*/
	orgId := "165" // This is the organisation ID
	/*
	   The App ID of the Workvivo app with unwiredotp.*.write permission.
	*/
	appWorkvivo := "3981" // APP ID from your Workvivo app with unwiredotp.*.write permission
	/*
	   The Domain associated with the Workvivo app
	*/
	appAud := "unwired.workvivo.red"
	/*
	   The API Gateway URL for Workvivo, sample code is using api.workvivo.red, for production environments you will need to change this to the appropriate URL for your Workvivo environment.

	   EU Production
	   api.workvivo.com
	   api.eu2.workvivo.com

	   US Production
	   api.workvivo.us
	   api.us2.workvivo.us

	   Middle East Production
	   api.workvivo.me

	*/
	apiURL := "https://api.workvivo.red/v1/unwired/users/otp"

	/*
	   This is the email address of the user you want a one-time passcode
	*/
	postData := `{"email": "test@nomail"}`

	/*
		this will be the customer’s host typically, e.g. acme.com
	*/
	appIssuer := "org1"

	/*
		this will always be 'app' as this API is always called as a Applicaton configured in Workvivo that the Partner Application is using to generate OTPs
	*/
	appSubject := "app"

	// Create JWT payload using struct for field order
	payload := JWTPayload{
		Aud:   appAud,
		AppID: appWorkvivo,
		Sub:   appSubject,
		Iss:   appIssuer,
		Nbf:   time.Now().Unix(),
		Iat:   time.Now().Unix(),
		Exp:   time.Now().Add(10 * time.Minute).Unix(),
		Jti:   randomHex(32),
		State: randomHex(32),
	}

	// Convert struct to map for JWT library
	payloadMap := map[string]interface{}{
		"aud":    payload.Aud,
		"app_id": payload.AppID,
		"sub":    payload.Sub,
		"iss":    payload.Iss,
		"nbf":    payload.Nbf,
		"iat":    payload.Iat,
		"exp":    payload.Exp,
		"jti":    payload.Jti,
		"state":  payload.State,
	}

	jwks := readJWKS()
	var kid string
	if keys, ok := jwks["keys"].([]interface{}); ok && len(keys) > 0 {
		if key, ok := keys[0].(map[string]interface{}); ok {
			if k, ok := key["kid"].(string); ok {
				kid = k
			}
		}
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims(payloadMap))
	token.Header["kid"] = kid
	jwtString, err := token.SignedString(privateKey)
	check(err)

	// Output JWT payload as prettified JSON
	payloadJSON, err := json.MarshalIndent(payload, "", "  ")
	check(err)
	fmt.Println("JWT payload (JSON):\n" + string(payloadJSON) + "\n")

	fmt.Println("JWT encoded:\n" + jwtString + "\n")
	fmt.Println("KeyID:\n" + kid + "\n")
	fmt.Println("Public Key (PEM):")
	printPublicKeyPEM(publicKey)
	fmt.Println("\nJWKS:")

	jwksJSON, _ := json.MarshalIndent(jwks, "", "  ")

	fmt.Println(string(jwksJSON))

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

func printPublicKeyPEM(pub *rsa.PublicKey) {
	pubBytes, err := x509.MarshalPKIXPublicKey(pub)
	check(err)
	pem.Encode(os.Stdout, &pem.Block{Type: "PUBLIC KEY", Bytes: pubBytes})
}

func loadKeyPair() (*rsa.PrivateKey, *rsa.PublicKey) {
	privKeyPath := "../Keys/private.pem"

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

func readJWKS() map[string]interface{} {
	jwks, err := os.ReadFile("../Keys/jwks.json")
	if err != nil {
		fmt.Println("Error reading JWKS file:", err)
		return map[string]interface{}{"keys": []interface{}{}}
	}
	var result map[string]interface{}
	err = json.Unmarshal(jwks, &result)
	if err != nil {
		fmt.Println("Error parsing JWKS JSON:", err)
		return map[string]interface{}{"keys": []interface{}{}}
	}
	return result
}
