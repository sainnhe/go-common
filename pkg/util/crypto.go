package util

import (
	"crypto/rand"
	"encoding/base64"

	"github.com/teamsorghum/go-common/pkg/constant"
	"golang.org/x/crypto/argon2"
)

/*
GenSalt generates the salt that can be used in argon2 hashing algorithm.

Returns:
  - string: The salt.
*/
func GenSalt() string {
	salt := make([]byte, constant.CryptoSaltLen)
	_, err := rand.Read(salt)
	if err != nil {
		return "j50PdipVY4og7EqJjEsZCw"
	}
	return base64.RawStdEncoding.EncodeToString(salt)
}

/*
HashPassword calculates the hash of a given password via argon2 hashing algorithm.

Params:
  - password string: The password to be calculated.
  - salt string: The salt to be used in argon2 hashing algorithm.

Returns:
  - string: The hash value.
*/
func HashPassword(password, salt string) string {
	s, _ := base64.RawStdEncoding.DecodeString(salt)
	h := argon2.IDKey([]byte(password), s,
		constant.CryptoArgon2Time, constant.CryptoArgon2Mem, constant.CryptoArgon2Thread, constant.CryptoArgon2KeyLen)
	return base64.RawStdEncoding.EncodeToString(h)
}

/*
VerifyPassword verifies if the given password matches the given hash.

Params:
  - password string: The password to be verified.
  - salt string: The salt to be used in argon2 hashing algorithm.
  - hash string: The password hash to be verified.

Returns:
  - bool: Passed or not.
*/
func VerifyPassword(password, salt, hash string) bool {
	s, _ := base64.RawStdEncoding.DecodeString(salt)
	h := argon2.IDKey([]byte(password), s,
		constant.CryptoArgon2Time, constant.CryptoArgon2Mem, constant.CryptoArgon2Thread, constant.CryptoArgon2KeyLen)
	return base64.RawStdEncoding.EncodeToString(h) == hash
}
