package num

import (
	"crypto/rand"
	"fmt"
	"math"
	"math/big"
	"strconv"
)

// Round will round the floating point number x to two decimal places.
func Round(x float64) float64 {
	return math.Round(x*100) / 100
}

func IntToBytes(i int) []byte {
	return []byte(strconv.Itoa(i))
}

func Float64ToString(x float64) string {
	return fmt.Sprintf("%.2f", Round(x))
}

func RandomInt(n int64) int64 {
	x, err := rand.Int(rand.Reader, big.NewInt(n))
	if err != nil {
		panic(err)
	}
	return x.Int64()
}
