package jwt

import (
	"encoding/base64"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var signKey = []byte("rabbit.RM")

func TestCustomClaims(t *testing.T) {
	type CustomClaims struct {
		jwt.RegisteredClaims
		Foo string `json:"foo"`
	}

	claims := CustomClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer: "RM",
		},
		Foo: "Bar",
	}
	fmt.Printf("claims:%v\n", claims)
	s, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(signKey)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(s)
}

func TestRegisteredClaims(t *testing.T) {
	claims := &jwt.RegisteredClaims{
		Issuer:    "RM",
		Subject:   "TEST",
		Audience:  []string{"TEST_1", "TEST_2"},
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(60 * time.Second)),
		NotBefore: jwt.NewNumericDate(time.Now()),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ID:        "1",
	}
	s, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(signKey)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(s)

}

func TestDecode(t *testing.T) {
	bb, err := base64.RawURLEncoding.DecodeString("eyJleHAiOjE3MTA0ODQ0MzMsIkZvbyI6IkJhciJ9")
	fmt.Println(string(bb), err)
}
