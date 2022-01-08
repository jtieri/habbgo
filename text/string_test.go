package text

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestContainsAllowedChars(t *testing.T) {
	allowedChars := "1234567890qwertyuiopasdfghjklzxcvbnm_-+=?!@:.,$"
	require.True(t, ContainsAllowedChars("1", allowedChars))
	require.False(t, ContainsAllowedChars(" lazar", allowedChars))
	require.False(t, ContainsAllowedChars("la zar", allowedChars))
	require.False(t, ContainsAllowedChars("lazar ", allowedChars))
	require.False(t, ContainsAllowedChars("~ghost", allowedChars))
	require.True(t, ContainsAllowedChars("____________", allowedChars))
	require.True(t, ContainsAllowedChars("0nem", allowedChars))
	require.True(t, ContainsAllowedChars("123456", allowedChars))
	require.False(t, ContainsAllowedChars("!---the winner---!", allowedChars))
	require.True(t, ContainsAllowedChars("____________", allowedChars))
}
