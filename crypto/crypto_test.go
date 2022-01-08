package crypto

import (
	"testing"
)

// TestLoginHash tests that the password hash/salt functions are working as intended
func TestLoginHash(t *testing.T) {
	var salt = GenerateRandomSalt(SALTSIZE)

	//loc, _ := time.LoadLocation("UTC")
	//now := time.Now().In(loc)
	//fmt.Println(now.Format("2006-01-02 15:04:05"))

	var hashedPassword = HashPassword("hello", salt)

	t.Log("Password Hash:", hashedPassword)
	t.Log("Salt:", salt)
	salty := string(salt)
	t.Log("Salt:", salty)
	salted := []byte(salty)
	t.Log("Salt:", salted)

	t.Log("Password Match:", PasswordsMatch(hashedPassword, "hello", salt))
}
