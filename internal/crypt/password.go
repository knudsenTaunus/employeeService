package crypt

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
)

type Encryptor struct {
	secret string
	bytes  []byte
}

func New(secret string) *Encryptor {
	return &Encryptor{
		secret: secret,
		bytes:  []byte{35, 46, 57, 24, 85, 35, 24, 74, 87, 35, 88, 98, 66, 32, 14, 05},
	}
}

func (e *Encryptor) Encrypt(text string) (string, error) {
	block, err := aes.NewCipher([]byte(e.secret))
	if err != nil {
		return "", err
	}
	plainText := []byte(text)
	cfb := cipher.NewCFBEncrypter(block, e.bytes)
	cipherText := make([]byte, len(plainText))
	cfb.XORKeyStream(cipherText, plainText)
	return base64.StdEncoding.EncodeToString(cipherText), nil
	//return Encode(cipherText), nil
}

func (e *Encryptor) Decrypt(text string) (string, error) {
	block, err := aes.NewCipher([]byte(e.secret))
	if err != nil {
		return "", err
	}

	cipherText, err := base64.StdEncoding.DecodeString(text)
	if err != nil {
		panic(err)
	}

	cfb := cipher.NewCFBDecrypter(block, e.bytes)
	plainText := make([]byte, len(cipherText))
	cfb.XORKeyStream(plainText, cipherText)
	return string(plainText), nil
}
