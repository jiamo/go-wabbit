package tests

import (
	"testing"
	"wabbit-go/model"
)

func TestNode(t *testing.T) {
	var env = map[string]int{}
	node := model.NewNodeInfo(env)
	println(node.Id)

}
