package model

import "testing"

func TestNode(t *testing.T) {
	var env = map[string]int{}
	node := NewNode(env)
	println(node.Id)
}
