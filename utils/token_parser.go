package utils

import (
	"fmt"

	"github.com/dgrijalva/jwt-go"
)

func ExtractClaims(tokenString string) (jwt.MapClaims, error) {
	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		return claims, nil
	} else {
		fmt.Println("Invalid jwt token")
		return nil, err
	}
}
