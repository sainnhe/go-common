// -------------------------------------------------------------------------------------------
// Copyright (c) Team Sorghum. All rights reserved.
// Licensed under the GPL v3 License. See LICENSE in the project root for license information.
// -------------------------------------------------------------------------------------------

package util_test

import (
	"fmt"
	"os"

	"github.com/teamsorghum/go-common/pkg/util"
)

// This example demonstrates how to generate a salt value and the hash value of a password and verify a incoming
// password using them.
//
// A typical usage is storing the generated salt and hash values into a database, and verifying the incoming password by
// looking at whether the password and the hash and salt values match.
func Example_crypto() {
	// Define your password
	password := "your_password"

	// Generate salt.
	salt, err := util.GenSalt()
	if err != nil {
		fmt.Printf("Error on generating salt: %s", err.Error())
		os.Exit(1)
	}

	// Calculate hash
	hash := util.HashPassword(password, salt)

	// Verify password
	result1 := util.VerifyPassword(password, salt, hash)

	// Verify another password
	result2 := util.VerifyPassword("wrong_password", salt, hash)

	fmt.Printf("Result1 = %t, result2 = %t\n", result1, result2)

	// Output: Result1 = true, result2 = false
}
