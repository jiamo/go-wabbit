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
		&model.PrintStatement{&model.Integer{"2"}},
		model.NewContext()))

	print(model.NodeAsSource(
		&model.Statements{
			[]model.Statement{
				&model.ConstDeclaration{model.Name{"pi"}, nil, &model.Float{"3.14159"}},

				&model.ConstDeclaration{model.Name{"tau"}, nil, &model.Mul{&model.Float{"2.0"}, &model.Name{"pi"}}},
				&model.VarDeclaration{model.Name{"radius"}, nil, &model.Float{"4.0"}},
				&model.VarDeclaration{model.Name{"perimeter"}, &model.NameType{"float"}, nil},
				&model.ExpressionAsStatement{&model.Assignment{&model.Name{"perimeter"}, &model.Mul{&model.Name{"tau"}, &model.Name{"radius"}}}},
				&model.PrintStatement{&model.Name{"perimeter"}},
			},
		},
		model.NewContext(),
	))

}
