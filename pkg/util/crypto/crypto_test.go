package cryptoutil_test

import (
	"encoding/base64"
	"testing"

	"github.com/teamsorghum/go-common/pkg/constant"
	cryptoutil "github.com/teamsorghum/go-common/pkg/util/crypto"
)

func TestGenSalt(t *testing.T) {
	t.Parallel()

	got := cryptoutil.GenSalt()
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
	salt := cryptoutil.GenSalt()
	got := cryptoutil.HashPassword(password, salt)

	if got == "" {
		t.Error("Generated hash is empty")
	}

	// The hash should be consistent given the same password and salt
	got2 := cryptoutil.HashPassword(password, salt)
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
	salt := cryptoutil.GenSalt()
	hash := cryptoutil.HashPassword(password, salt)

	// Test with the correct password
	if !cryptoutil.VerifyPassword(password, salt, hash) {
		t.Error("Correct password did not pass verification")
	}

	// Test with an incorrect password
	if cryptoutil.VerifyPassword("wrongpassword", salt, hash) {
		t.Error("Incorrect password passed verification")
	}
}
