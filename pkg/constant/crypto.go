// -------------------------------------------------------------------------------------------
// Copyright (c) Team Sorghum. All rights reserved.
// Licensed under the GPL v3 License. See LICENSE in the project root for license information.
// -------------------------------------------------------------------------------------------

package constant

const (
	// CryptoSaltLen is the length of salt.
	CryptoSaltLen = 16
	// CryptoDefaultSalt is the default salt value.
	CryptoDefaultSalt = "j50PdipVY4og7EqJjEsZCw"
	// CryptoArgon2Time is the time required by argon2.
	CryptoArgon2Time = 1
	// CryptoArgon2Mem is the memory required by argon2.
	CryptoArgon2Mem = 64 * 1024
	// CryptoArgon2Thread is the thread required by argon2.
	CryptoArgon2Thread = 4
	// CryptoArgon2KeyLen is the key length required by argon2.
	CryptoArgon2KeyLen = 32
)
