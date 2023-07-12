package main

import "wabbit-go/model"

func main() {
	modelA := &model.BinOpWithOp{
		"+",
		&model.Integer{"2"},
		&model.BinOpWithOp{
			"*",
			&model.Integer{"2"},
			&model.Integer{"4"}}}

	println(model.NodeAsSource(modelA, model.NewContext()))
	println(model.NodeAsSource(
		&model.PrintStatement{&model.Integer{"2"}}, model.NewContext()))
}
