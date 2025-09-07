package utils

import (
	"fmt"
	"math/rand"
	"time"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
)

// GenerateOTP returns a 4-digit numeric OTP
func GenerateOTP() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%04d", rand.Intn(10000))
}




// Encrypt OTP
func EncryptOTP(otp string, key string) (string, error) {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}
	plaintext := []byte(otp)
	cfb := cipher.NewCFBEncrypter(block, []byte(key)[:block.BlockSize()])
	ciphertext := make([]byte, len(plaintext))
	cfb.XORKeyStream(ciphertext, plaintext)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// Decrypt OTP
func DecryptOTP(encrypted string, key string) (string, error) {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}
	ciphertext, _ := base64.StdEncoding.DecodeString(encrypted)
	cfb := cipher.NewCFBDecrypter(block, []byte(key)[:block.BlockSize()])
	plaintext := make([]byte, len(ciphertext))
	cfb.XORKeyStream(plaintext, ciphertext)
	return string(plaintext), nil
}
