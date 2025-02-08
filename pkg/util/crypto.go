package util

import (
	"crypto/rand"
	"encoding/base64"

	"github.com/teamsorghum/go-common/pkg/constant"
	"golang.org/x/crypto/argon2"
)

// GenSalt generates the salt that can be used in argon2 hashing algorithm.
func GenSalt() string {
	salt := make([]byte, constant.CryptoSaltLen)
	_, err := rand.Read(salt)
	if err != nil {
		return constant.CryptoDefaultSalt
	}
	return base64.RawStdEncoding.EncodeToString(salt)
}

// HashPassword calculates the hash of a given password via argon2 hashing algorithm.
func HashPassword(password, salt string) string {
	s, _ := base64.RawStdEncoding.DecodeString(salt)
	h := argon2.IDKey([]byte(password), s,
		constant.CryptoArgon2Time, constant.CryptoArgon2Mem, constant.CryptoArgon2Thread, constant.CryptoArgon2KeyLen)
	return base64.RawStdEncoding.EncodeToString(h)
}

// VerifyPassword verifies if the given password matches the given hash.
func VerifyPassword(password, salt, hash string) bool {
	s, _ := base64.RawStdEncoding.DecodeString(salt)
	h := argon2.IDKey([]byte(password), s,
		constant.CryptoArgon2Time, constant.CryptoArgon2Mem, constant.CryptoArgon2Thread, constant.CryptoArgon2KeyLen)
	return base64.RawStdEncoding.EncodeToString(h) == hash
}
