package collections

import (
	"testing"
)

type MockValue struct {
}

func TestGenericCache_SetGetContains(t *testing.T) {
	m := make(map[int]*MockValue)
	c := NewCache(m)

	const key = 100
	value := &MockValue{}
	c.Set(key, value)

	v, ok := c.Get(key)
	if ok == false {
		t.Errorf("item with key %d should have been present in cache", key)
	}
	if v == nil {
		t.Error("item retrieval should not have returned a nil value")
	}

	contains := c.Has(key)
	if contains == false {
		t.Error("call to Has should have returned true")
	}
}
