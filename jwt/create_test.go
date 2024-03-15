package jwt

import (
	"encoding/base64"
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func TestNewClaims(t *testing.T) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user": "admin",
		"exp":  jwt.NewNumericDate(time.Now()),
		"iss":  "rabbit.RM",
	})
	s, err := token.SignedString([]byte("rabbit.RM"))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(s)
	split := strings.Split(s, ".")
	for _, item := range split {
		ds, err := base64.RawURLEncoding.DecodeString(item)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(ds))
	}

}

func TestCreate(t *testing.T) {
	key := ([]byte)("rabbit.RM")
	token := jwt.New(jwt.SigningMethodHS256)
	s, err := token.SignedString(key)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(s)

	split := strings.Split(s, ".")
	for _, item := range split {
		ds, err := base64.RawURLEncoding.DecodeString(item)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(ds))
	}
}
