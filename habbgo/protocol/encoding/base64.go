/*
base64 contains an implementation of the FUSE-Base64 tetrasexagesimal numeric system used in the FUSE v0.2.0 protocol.
It typically uses two ASCII characters between decimal indexes 64 (@) and 127 (DEL control character) to produce a
two-character representation of a number between 0 and 4095.

This implementation is a Golang port of the examples found on Puomi's wiki page for Base64.
*/
package encoding

import "math"

// EncodeB64 takes an integer, encodes it in FUSE-Base64 & returns a slice of, length number of, bytes.
func EncodeB64(i int, length int) []byte {
	bytes := make([]byte, length)
	for j := 1; j <= length; j++ {
		k := uint((length - j) * 6)
		bytes[j-1] = byte(0x40 + ((i >> k) & 0x3f))
	}
	return bytes
}

// DecodeB64 take a slice of bytes, decodes it from FUSE-Base64 & returns the decoded bytes as an integer.
func DecodeB64(bytes []byte) int {
	decodedVal := 0
	for i, j := len(bytes) - 1, 0; i >= 0; i-- {
		x := int(bytes[i] - 0x40)
		if j > 0 {
			x *= int(math.Pow(64.0, float64(j)))
		}

		decodedVal += x
		j--
	}
	return decodedVal
}
