package lib

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"fmt"
)

type Signature struct {
	logger Logger
	env    Env
}

func NewSignature(logger Logger, env Env) Signature {
	return Signature{
		logger: logger,
		env:    env,
	}
}

var key []byte
var nonce []byte

func (s Signature) Encrypt(encryptedString string) string {

	key, _ = hex.DecodeString(s.env.ClientSecretKey)
	nonce = make([]byte, 12)

	plaintext := []byte(encryptedString)

	// Create a new Cipher Block from the key
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	// Create a new GCM
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		panic(err)
	}

	ciphertext := aesGCM.Seal(nil, nonce, plaintext, nil)
	return fmt.Sprintf("%x", ciphertext)
}

func (s Signature) Decrypt(encryptedString string) string {
	enc, _ := hex.DecodeString(encryptedString)

	// Create a new Cipher Block from the key
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	// Create a new GCM
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		panic(err)
	}

	// Decrypt the data
	plaintext, err := aesGCM.Open(nil, nonce, enc, nil)
	if err != nil {
		panic(err)
	}

	return string(plaintext)
}
