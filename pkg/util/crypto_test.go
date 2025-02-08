package util_test

import (
	"encoding/base64"
	"testing"

	"github.com/teamsorghum/go-common/pkg/constant"
	"github.com/teamsorghum/go-common/pkg/util"
)

func TestGenSalt(t *testing.T) {
	t.Parallel()

	got := util.GenSalt()
	if got == "" {
		t.Error("Generated salt is empty")
	}

	// Check if the length of the salt matches the expected length
	wantSaltLen := base64.RawStdEncoding.EncodedLen(constant.CryptoSaltLen)
	if len(got) != wantSaltLen {
		t.Errorf("Want salt length %d, got %d", wantSaltLen, len(got))
	}
}

func TestHashPassword(t *testing.T) {
	t.Parallel()

	password := "password123"
	salt := util.GenSalt()
	got := util.HashPassword(password, salt)

	if got == "" {
		t.Error("Generated hash is empty")
	}

	// The hash should be consistent given the same password and salt
	got2 := util.HashPassword(password, salt)
	if got != got2 {
		t.Error("Hashes do not match for the same password and salt")
	}

	// Check if the length of the hash matches the expected length
	wantHashLen := base64.RawStdEncoding.EncodedLen(int(constant.CryptoArgon2KeyLen))
	if len(got) != wantHashLen {
		t.Errorf("Want hash length %d, got %d", wantHashLen, len(got))
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
