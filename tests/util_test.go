package tests

import (
	"fmt"
	"testing"
	"wabbit-go/common"
)

func TestChainMap(t *testing.T) {
	m1 := common.NewChainMap()
	m1.SetValue("key1", "value1")
	m1.SetValue("key2", "value2")

	m2 := m1.NewChild()
	m2.SetValue("key3", "value3")
	m2.SetValue("key4", "value4")

	value, ok := m1.GetValue("key3")
	fmt.Printf("key3: %v, ok: %v\n", value, ok)

	value, ok = m2.GetValue("key1")
	fmt.Printf("key1: %v, ok: %v\n", value, ok)

	value, ok = m2.GetValue("key3")
	fmt.Printf("key3: %v, ok: %v\n", value, ok)

	value, ok = m2.GetValue("key5")
	fmt.Printf("key5: %v, ok: %v\n", value, ok)
}
