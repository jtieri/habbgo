package date

import (
	"testing"
)

func TestGetCurrentDate(t *testing.T) {
	t.Log(GetCurrentDate())
	t.Log(GetCurrentDateTime())
}
