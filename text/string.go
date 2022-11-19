package text

import (
	"crypto/rand"
	"math/big"
	"strconv"
	"strings"
)

// AllowedCharacters are the characters that are permitted to be used in usernames and password.
const AllowedCharacters = "1234567890qwertyuiopasdfghjklzxcvbnm_-+=?!@:.,$" // TODO possibly make this configurable in the future.

// Filter will replace special characters found in the specified string with an empty string.
func Filter(s string) string {
	output := strings.Replace(s, string(rune(1)), "", -1)
	output = strings.Replace(output, string(rune(2)), "", -1)
	output = strings.Replace(output, string(rune(9)), "", -1)
	output = strings.Replace(output, string(rune(10)), "", -1)
	output = strings.Replace(output, string(rune(12)), "", -1)
	output = strings.Replace(output, string(rune(13)), "", -1) // remove newlines
	return output
}

// ContainsAllowedChars returns false if the toTest string contains a character that is not present in AllowedCharacters and true otherwise.
func ContainsAllowedChars(toTest, allowedChars string) bool {
	for _, v := range toTest {
		if strings.Contains(allowedChars, string(v)) {
			continue
		}
		return false
	}

	return true
}

// ContainsNumber returns true if the specified string contains a numerical character in it and false otherwise.
func ContainsNumber(text string) bool {
	for _, l := range text {
		if _, err := strconv.Atoi(string(l)); err == nil {
			return true
		}
	}

	return false
}

// Substr returns a sub string of the string input from the start index to length-1.
func Substr(input string, start int, length int) string {
	asRunes := []rune(input)

	if start >= len(asRunes) {
		return ""
	}

	if start+length > len(asRunes) {
		length = len(asRunes) - start
	}

	return string(asRunes[start : start+length])
}

var chars = []rune("abcdefghijklmnopqrstuvwxyz")

// RandomString returns a string of random characters for the specified length.
func RandomString(length int) string {
	var b strings.Builder
	for i := 0; i < length; i++ {
		i, _ := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
		b.WriteRune(chars[i.Int64()])
	}
	return b.String()
}
