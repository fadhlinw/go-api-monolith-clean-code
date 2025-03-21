package utils

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

func GenerateRandomString(n int) string {
	var letterRunes = []rune("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz")

	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

// Generate random 6 digit number string
func GenerateRandomNumberString() string {
	rand.Seed(time.Now().UnixNano())
	min := 100000
	max := 999999
	return fmt.Sprintf("%d", rand.Intn(max-min+1)+min)
}

func ReplaceComaWithSpace(key string) string {
	return strings.Replace(key, ",", " ", -1)
}
