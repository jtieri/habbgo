package crypto

import (
	"crypto/rand"
	"crypto/sha512"
	"encoding/base64"
)

const SALTSIZE = 16 // default password salt size in bytes

// GenerateRandomSalt creates a new slice of bytes and generates a random salt using the cryptographically secure CSPRNG
func GenerateRandomSalt(saltSize int) []byte {
	salt := make([]byte, saltSize)
	if _, err := rand.Read(salt[:]); err != nil {
		panic(err) // TODO handle error gracefully
	}
	return salt
}

// HashPassword will concatenate the provided password and salt, then use SHA-512 to hash the password
// before returning it as a base64 encoded string
func HashPassword(password string, salt []byte) string {
	var sha512Hasher = sha512.New() // Create a sha-512 hasher

	passwordBytes := []byte(password)
	passwordBytes = append(passwordBytes, salt...)

	sha512Hasher.Write(passwordBytes)
	hashedPasswordBytes := sha512Hasher.Sum(nil) // Get the SHA-512 hashed password

	// Base64 encode the hashed password
	return base64.URLEncoding.EncodeToString(hashedPasswordBytes)
}

// PasswordsMatch ensures an un-hashed string matches a hashed string once it is combined with a salt & hashed itself
func PasswordsMatch(hashedPassword, unhashedPass string, salt []byte) bool {
	return hashedPassword == HashPassword(unhashedPass, salt)
}
