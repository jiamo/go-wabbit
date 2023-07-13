package main

import "wabbit-go/model"

func main() {
	modelA := &model.BinOpWithOp{
		"+",
		&model.Integer{2},
		&model.BinOpWithOp{
			"*",
			&model.Integer{2},
			&model.Integer{4}}}

	println(model.NodeAsSource(modelA, model.NewContext()))
	println(model.NodeAsSource(
		&model.PrintStatement{&model.Integer{2}},
		model.NewContext()))

	println(model.NodeAsSource(
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

	println(model.NodeAsSource(
		&model.PrintStatement{&model.Eq{&model.Integer{2}, &model.Integer{2}}},
		model.NewContext()))

	println(model.NodeAsSource(
		&model.PrintStatement{
			&model.LogOr{&model.NameBool{"true"},
				&model.Grouping{
					&model.Eq{
						&model.Div{&model.Integer{1}, &model.Integer{2}},
						&model.Integer{2},
					}}}},
		model.NewContext()))

	println(model.NodeAsSource(
		&model.Statements{
			[]model.Statement{
				&model.ConstDeclaration{model.Name{"a"}, &model.NameType{"int"}, &model.Integer{2}},
				&model.ConstDeclaration{model.Name{"b"}, &model.NameType{"int"}, &model.Integer{3}},
				&model.IfStatement{
					&model.Lt{&model.Name{"a"}, &model.Name{"b"}},
					model.Statements{[]model.Statement{
						&model.ExpressionAsStatement{&model.Assignment{&model.Name{"minvar"}, &model.Name{"a"}}},
					}},
					&model.Statements{[]model.Statement{
						&model.ExpressionAsStatement{&model.Assignment{&model.Name{"minvar"}, &model.Name{"b"}}},
					}},
				},
			},
		},
		model.NewContext(),
	))

	p3 := &model.Statements{
		[]model.Statement{
			&model.VarDeclaration{model.Name{"n"}, &model.IntegerType{}, &model.Integer{'1'}},
			&model.WhileStatement{
				&model.NameBool{"true"},
				model.Statements{[]model.Statement{
					&model.IfStatement{
						&model.Eq{&model.Name{"n"}, &model.Integer{2}},
						model.Statements{
							[]model.Statement{
								&model.PrintStatement{&model.Name{"n"}},
								&model.BreakStatement{},
							}},
						&model.Statements{
							[]model.Statement{
								&model.ExpressionAsStatement{&model.Assignment{&model.Name{"n"}, &model.Add{&model.Name{"n"}, &model.Integer{'1'}}}},
								&model.ContinueStatement{},
							}},
					},
					&model.ExpressionAsStatement{&model.Assignment{&model.Name{"n"}, &model.Sub{&model.Name{"n"}, &model.Integer{1}}}}},
				}},
			&model.PrintStatement{&model.Name{"n"}}},
	}
	println(model.NodeAsSource(p3, model.NewContext()))
}
