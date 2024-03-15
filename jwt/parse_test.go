package jwt

import (
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func TestCustomClaimsValidation(t *testing.T) {
	type CustomClaims struct {
		Foo string `json:"foo"`
		jwt.RegisteredClaims
	}
	tokenStr := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJSTSIsImZvbyI6IkJhciJ9.AVZbvGL3T9STHSpOjdoV6exEjdCCbU4B42SmXlpI-Ys"
	token, err := jwt.NewParser(jwt.WithLeeway(5*time.Minute)).ParseWithClaims(tokenStr, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return signKey, nil
	})
	if err != nil {
		log.Fatal(err)
	} else if v, ok := token.Claims.(*CustomClaims); ok {
		fmt.Println(v.Foo, v.Issuer)
	} else {
		log.Fatal("unknown claims type,cannot proceed")
	}

}

func TestParseCustomClaims(t *testing.T) {
	type CustomClaims struct {
		Foo string `json:"foo"`
		jwt.RegisteredClaims
	}
	tokenStr := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJSTSIsImZvbyI6IkJhciJ9.AVZbvGL3T9STHSpOjdoV6exEjdCCbU4B42SmXlpI-Ys"
	token, err := jwt.ParseWithClaims(tokenStr, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return signKey, nil
	})
	if err != nil {
		log.Fatal(err)
	} else if v, ok := token.Claims.(*CustomClaims); ok {
		fmt.Println(v.Foo, v.Issuer)
	} else {
		log.Fatal("unknown claims type,cannot proceed")
	}
}
