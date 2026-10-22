// Package otp creates and checks the one-time codes used for login nd phone verification
package otp

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"math/big"
)

// Generate returns a random 6-digit code such as "042917"
func Generate() (string, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(1_000_000))
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%06d", n.Int64()), nil
}

// Hash returns the HMAC-SHA256 of code, keyed with secret.
func Hash(secret []byte, code string) []byte {
	mac := hmac.New(sha256.New, secret)
	mac.Write([]byte(code))
	return mac.Sum(nil)
}

// Matches reports whether code hashes to want, comparing in constant time.
func Matches(secret []byte, code string, want []byte) bool {
	return hmac.Equal(Hash(secret, code), want)
}
