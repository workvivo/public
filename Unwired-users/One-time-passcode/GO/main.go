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

	// Create JWT
	payload := jwt.MapClaims{
		"jti":         randomHex(32),
		"iss":         appIssuer,
		"sub":         appSubject,
		"workvivo_id": appWorkvivo,
		"aud":         appAud,
		"iat":         time.Now().Unix(),
		"nbf":         time.Now().Unix(),
		"exp":         time.Now().Add(10 * time.Minute).Unix(),
		"state":       randomHex(32),
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

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, payload)
	token.Header["kid"] = kid
	jwtString, err := token.SignedString(privateKey)
	check(err)

	fmt.Println("Token (JWT):\n" + jwtString + "\n")
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
