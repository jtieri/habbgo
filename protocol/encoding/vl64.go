/*
vl64 contains an implementation of the FUSE mixed radix numeric encoding used by Sulake in HH.
This implementation is a Golang port of the examples found on Puomi's wiki page for VL64.
*/
package encoding

import "math"

// DecodeVl64 returns a single number from the Vl64 encoded input.
// Any characters after the length indicated by first char of input will be discarded.
func DecodeVl64(input []byte) int {
	length := length(input[0])
	total := int(input[0]) % 4 // Base4 value

	// Increment all Base64 symbols to the total
	for inc := 0; inc < length; inc++ {
		total += (int(input[inc]) - 64) * int(math.Pow(64, float64(inc))/16)
	}

	if int(input[0]%8) < 4 {
		return total // Base4 positive
	}

	return -total // Base4 negative
}

// EncodeVl64 returns a slice of bytes capable of storing all increments.
// Removes all @ symbols (padding) after last non-@ before returning.
func EncodeVl64(input int) []byte {
	vl64 := make([]byte, 6)              // 32-bit integer causes VL64 to have max length of 6.
	num := int(math.Abs(float64(input))) // Operate on normalized, positive integer
	length := 1                          // Length indicator, updated during encode

	var indicator int
	if input < 0 {
		indicator = 'D' // positive(+64)
	} else {
		indicator = '@' // negative(+68)
	}
	vl64[0] = byte(num%4 + indicator) // Base4 char, positive(+64)/negative(+68) indicator
	num /= 4                          // Base4 processed, prepare for remaining Base64 symbols

	for i := 1; i < 6; i++ {
		vl64[i] = byte(num%64 + 64) // Base64
		num /= 64

		if vl64[i] != 64 {
			length = i + 1 // @ = padding / zero symbol
		}

	}

	vl64[0] = byte(int(vl64[0]) + length*8) // Base4 char shifted to indicate total length
	return vl64[:length]                    // Last padding symbols trimmed out
}

// length returns the total length of the mixed radix number.
func length(firstChar byte) int {
	return (int(firstChar) - 64) / 8
}
