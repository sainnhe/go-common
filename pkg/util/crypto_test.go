package util_test

import (
	"encoding/base64"
	"testing"

	"github.com/teamsorghum/go-common/pkg/constant"
	"github.com/teamsorghum/go-common/pkg/util"
)

func TestGenSalt(t *testing.T) {
	t.Parallel()

	salt := util.GenSalt()
	if salt == "" {
		t.Error("Generated salt is empty")
	}

	// Check if the length of the salt matches the expected length
	expectedSaltLen := base64.RawStdEncoding.EncodedLen(constant.CryptoSaltLen)
	if len(salt) != expectedSaltLen {
		t.Errorf("Expected salt length %d, got %d", expectedSaltLen, len(salt))
	}
}

func TestHashPassword(t *testing.T) {
	t.Parallel()

	password := "password123"
	salt := util.GenSalt()
	hash := util.HashPassword(password, salt)

	if hash == "" {
		t.Error("Generated hash is empty")
	}

	// The hash should be consistent given the same password and salt
	hash2 := util.HashPassword(password, salt)
	if hash != hash2 {
		t.Error("Hashes do not match for the same password and salt")
	}

	// Check if the length of the hash matches the expected length
	expectedHashLen := base64.RawStdEncoding.EncodedLen(int(constant.CryptoArgon2KeyLen))
	if len(hash) != expectedHashLen {
		t.Errorf("Expected hash length %d, got %d", expectedHashLen, len(hash))
	}
}

func TestVerifyPassword(t *testing.T) {
	t.Parallel()

	password := "password123"
	salt := util.GenSalt()
	hash := util.HashPassword(password, salt)

	// Test with the correct password
	if !util.VerifyPassword(password, salt, hash) {
		t.Error("Correct password did not pass verification")
	}

	// Test with an incorrect password
	if util.VerifyPassword("wrongpassword", salt, hash) {
		t.Error("Incorrect password passed verification")
	}
}
