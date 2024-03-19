package jwt

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/golang-jwt/jwt/v5"
)

func TestRsaCreateToken(t *testing.T) {
	bytes, err := os.ReadFile("./testdata/rsa.private.pem")
	if err != nil {
		log.Fatal(err)
	}
	pem, err := jwt.ParseRSAPrivateKeyFromPEM(bytes)
	if err != nil {
		log.Fatal(err)
	}
	claims := &jwt.RegisteredClaims{Issuer: "Rabbit.RM"}
	str, err := jwt.NewWithClaims(jwt.SigningMethodRS512, claims).SignedString(pem)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(str)
}
