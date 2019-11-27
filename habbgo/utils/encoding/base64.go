/*
base64 contains an implementation of the FUSE-Base64 tetrasexagesimal numeric system used in the FUSE v0.2.0 protocol.
It typically uses two ASCII characters between decimal indexes 64 (@) and 127 (DEL control character) to produce a
two-character representation of a number between 0 and 4095.

This implementation is a Golang port of the examples from Puomi's wiki for Base64.
*/
package encoding

// EncodeB64 takes an integer, encodes it in FUSE-Base64 & returns a slice of bytes that should contain two char's.
func EncodeB64(i int) []byte {
	return []byte{byte(i/64 + 64), byte(i%64 + 64)}
}

// DecodeB64 take a slice of bytes, decodes it from FUSE-Base64 & returns the decoded bytes as an integer.
func DecodeB64(bytes []byte) int {
	return 64 * (int(bytes[0]%64) + int(bytes[1]%64))
}
