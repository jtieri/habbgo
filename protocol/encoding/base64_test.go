package encoding

import "testing"

func Test_Base64Encode(t *testing.T) {

}

func Test_Base64Decode(t *testing.T) {
	input := "A"
	inputBz := []byte(input)
	t.Log(DecodeB64(inputBz))

}
